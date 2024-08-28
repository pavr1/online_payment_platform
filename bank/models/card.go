package models

type Card struct {
	ID         string    `json:"id"`
	CardNumber string    `json:"card_number"`
	HolderName string    `json:"holder_name"`
	ExpDate    string    `json:"exp_date"`
	CVV        string    `json:"cvv"`
	Account    *Account  `json:"account"`
	Customer   *Customer `json:"customer"`
}

func (c *Card) GetAmount() float64 {
	return c.Account.Amount
}

func (c *Card) SetAmount(amount float64) {
	c.Account.Amount = amount
}

func (c *Card) GetCustomerName() string {
	return c.Customer.FirstName + " " + c.Customer.LastName
}
