package models

type Card struct {
	ID         string    `json:"id"`
	CardNumber string    `json:"card_number"`
	HolderName string    `json:"holder_name"`
	ExpDate    string    `json:"exp_date"`
	CVV        string    `json:"cvv"`
	account    *Account  `json:"account"`
	customer   *Customer `json:"customer"`
}

func (c *Card) GetAmount() float64 {
	return c.account.Amount
}

func (c *Card) SetAmount(amount float64) {
	c.account.Amount = amount
}

func (c *Card) GetCustomerName() string {
	return c.customer.FirstName + " " + c.customer.LastName
}
