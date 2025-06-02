package account

import "github.com/google/uuid"

type Transaction struct {
	ID          string `json:"id"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
	Category    string `json:"category"`
}

func NewTransaction(date, description string, amount int64, category string) Transaction {
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

func FindIndexByID(id string, transactionList []Transaction) int {
	for i, t := range transactionList {
		if t.ID == id {
			return i
		}
	}
	return -1
}
