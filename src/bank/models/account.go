package models

import (
	"github.com/pavr1/online_payment_platform/bank/models/base"
)

type Account struct {
	base.BaseModel
	ID            string
	AccountNumber string
	Amount        float32
}
