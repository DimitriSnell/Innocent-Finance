package account

import "github.com/google/uuid"

type Transaction struct {
	ID          string `json:"id"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
	Category    string `json:"category"`
	DonatorID   string `json:"donatorid"`
}

func NewTransaction(date, description string, amount int64, category string, DID string) Transaction {
	id := GenerateID()
	return Transaction{
		ID:          id,
		Date:        date,
		Description: description,
		Amount:      amount,
		Category:    category,
		DonatorID:   DID,
	}
}

func GenerateID() string {
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
