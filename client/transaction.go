package main

import "github.com/google/uuid"

type Transaction struct {
	ID          string `json:"id"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
	Category    string `json:"category"`
}

func newTransaction(date, description string, amount int64, category string) Transaction {
	id := generateID()
	return Transaction{
		ID:          id,
		Date:        date,
		Description: description,
		Amount:      amount,
		Category:    category,
	}
}

func generateID() string {
	return uuid.New().String()
}
