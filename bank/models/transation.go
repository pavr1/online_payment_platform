package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID        string             `json:"id"`
	Date      primitive.DateTime `json:"date"`
	Amount    float64            `json:"amount"`
	FromCard  string             `json:"from_card"`
	ToAccount string             `json:"to_account"`
	Detail    string             `json:"details"`
	Status    string             `json:"status"`
}
