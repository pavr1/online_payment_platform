package models

import (
	"github.com/pavr1/online_payment_platform/bank/models/base"
)

type Customer struct {
	base.BaseModel
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
