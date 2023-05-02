package models

import "time"

type Data struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int
}
type Order struct {
	ID        int64
	Price     float64
	Amount    float64
	Total     float64
	OrderType string
}

type OrderBookCup struct {
	Bids []Order
	Asks []Order
}
