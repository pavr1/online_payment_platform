package models

type Account struct {
	ID            string  `json:"id"`
	AccountNumber string  `json:"account_number"`
	Amount        float64 `json:"amount"`
}
