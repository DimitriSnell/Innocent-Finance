package main

type Transaction struct {
	ID          string `json:"id"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
	Category    string `json:"category"`
}

type Budget struct {
	MonthlyLimit int64 `json:"monthly_limit"`
	TotalBudget  int64 `json:"total_budget"`
}

type Data struct {
	Transactions []Transaction `json:"transactions"`
	Budget       Budget        `json:"budget"`
}

type Account struct {
	data Data
}

func newAccount(d Data) *Account {
	a := Account{d}
	return &a
}
