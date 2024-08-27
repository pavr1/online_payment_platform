package models

import (
	"time"

	"github.com/pavr1/online_payment_platform/bank/models/base"
)

type Transaction struct {
	base.BaseModel
	ID          string
	Date        time.Time
	Amount      float32
	FromAccount string
	ToAccount   string
	Detail      string
}
