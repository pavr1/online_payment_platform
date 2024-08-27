package models

type Card struct {
	CardNumber string
	HolderName string
	ExpDate    string
	CVV        string
	account    *Account
}
