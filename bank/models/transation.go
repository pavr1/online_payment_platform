package models

import (
	"time"
)

type Transaction struct {
	ID        string    `json:"id"`
	Date      time.Time `json:"date"`
	Amount    float64   `json:"amount"`
	FromCard  string    `json:"from_card"`
	ToAccount string    `json:"to_account"`
	Detail    string    `json:"details"`
}
