package models

import "time"

type Transaction struct {
	ID          string
	Date        time.Time
	Amount      float32
	FromAccount string
	ToAccount   string
	Detail      string
}
