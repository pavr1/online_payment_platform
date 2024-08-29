package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID        string             `json:"id"`
	Date      primitive.DateTime `json:"date"`
	Amount    float64            `json:"amount"`
	FromCard  string             `json:"fromcard"`
	ToAccount string             `json:"toaccount"`
	Detail    string             `json:"detail"`
	Status    string             `json:"status"`
}
