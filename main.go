package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	http.HandleFunc("/", serveHome)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})

	sr := NewStockRepo()
	stocks := sr.All()

	mu := sync.Mutex{}

	for _, stock := range stocks {
		go func(s string) {
			for {
				price := sr.CurrentPrice(s)
				fmt.Println(s, price)

				for client, subscribed := range clients {
					if _, ok := subscribed[s]; ok {
						mu.Lock()
						wr, err := client.NextWriter(websocket.TextMessage)
						if err != nil {
							log.Println("[err] writing to client", err)
						}

						res, err := json.Marshal(StockResponse{
							Stock: s,
							Price: price,
						})

						wr.Write(res)
						mu.Unlock()
					}
				}
				time.Sleep(2 * time.Second)
			}
		}(stock)
	}

	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
