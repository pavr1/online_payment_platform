package models

import (
	"github.com/pavr1/online_payment_platform/bank/models/base"
)

type Account struct {
	base.BaseModel
	ID            string  `json:"id"`
	AccountNumber string  `json:"account_number"`
	Amount        float32 `json:"amount"`
}
