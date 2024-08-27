package models

import (
	"time"
)

type Transaction struct {
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	Amount      float64   `json:"amount"`
	FromAccount string    `json:"from_account"`
	ToAccount   string    `json:"to_account"`
	Detail      string    `json:"details"`
}
