package main

func FindIndexByID(id string, transactionList []Transaction) int {
	for i, t := range transactionList {
		if t.ID == id {
			return i
		}
	}
	return -1
}
