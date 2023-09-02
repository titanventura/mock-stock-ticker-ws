package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	subscribe   = "SUBSCRIBE"
	unSubscribe = "UN_SUBSCRIBE"
)

var clients = make(map[*websocket.Conn]map[string]bool)

type Message struct {
	Type  string `json:"type"`
	Stock string `json:"stock"`
}

func handleClient(c *websocket.Conn) {
	defer func() {
		delete(clients, c)
		c.Close()
	}()
	clients[c] = make(map[string]bool)

	for {
		var m Message
		err := c.ReadJSON(&m)
		if err != nil {
			log.Println("[err] reading JSON from client", err)
		}

		if m.Type == subscribe {
			clients[c][m.Stock] = true
		} else if m.Type == unSubscribe {
			delete(clients[c], m.Stock)
		}
	}
}

type StockResponse struct {
	Stock string  `json:"stock"`
	Price float64 `json:"price"`
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go handleClient(conn)
}
