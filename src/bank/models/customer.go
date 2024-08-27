package models

import (
	"github.com/pavr1/online_payment_platform/bank/models/base"
)

type Customer struct {
	base.BaseModel
	ID        string
	FirstName string
	LastName  string
	Email     string
}
