package main

import (
	"math/rand"
	"time"
)

type StockRepo interface {
	All() []string
	CurrentPrice(id string) float64
}

type stockRepo struct {
	listOfStocks []string
}

func NewStockRepo() StockRepo {
	return &stockRepo{
		listOfStocks: []string{"AAPL", "GOOGL", "TWKS", "MSFT"},
	}
}

func (s *stockRepo) All() []string {
	return s.listOfStocks
}
func (s *stockRepo) CurrentPrice(id string) float64 {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return r.Float64()
}
