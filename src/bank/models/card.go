package models

import (
	"github.com/pavr1/online_payment_platform/bank/models/base"
)

type Card struct {
	base.BaseModel
	ID         string
	CardNumber string
	HolderName string
	ExpDate    string
	CVV        string
	account    *Account
}
