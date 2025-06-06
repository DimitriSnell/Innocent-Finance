package main

type Transaction struct {
	ID          string `json:"id"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
	Category    string `json:"category"`
	DonatorID   string `json:"donatorid"`
}

type Budget struct {
	MonthlyLimit int64 `json:"monthly_limit"`
	TotalBudget  int64 `json:"total_budget"`
}

type Data struct {
	Donators     []Donator     `json:"donators"`
	Transactions []Transaction `json:"transactions"`
	Budget       Budget        `json:"budget"`
	SyncToken    int64         `json:"sync_token"`
}

type Changes struct {
	AddedTransactions    []Transaction `json:"addedTransactions,omitempty"`
	DeletedTransactions  []Transaction `json:"deletedTransactions,omitempty"`
	ReplacedTransactions []Transaction `json:"replacedTransactions,omitempty"`
	AddedDonators        []Donator     `json:"addedDonator,omitempty"`
}

type Account struct {
	data    Data
	changes Changes
}

func newAccount(d Data) *Account {
	a := Account{d, Changes{}}
	return &a
}
