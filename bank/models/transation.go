package models

import (
	"time"

	"github.com/pavr1/online_payment_platform/bank/models/base"
)

type Transaction struct {
	base.BaseModel
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	Amount      float32   `json:"amount"`
	FromAccount string    `json:"from_account"`
	ToAccount   string    `json:"to_account"`
	Detail      string    `json:"details"`
}
