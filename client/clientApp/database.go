package clientapp

import (
	"client/account"
	"database/sql"
	"fmt"
	"strings"
)

type TransactionFilterInfo struct {
	ID          string
	Date        string
	Description string
	Amount      int64
	Category    string
}

func (c *Client) initDatabase() error {
	db, err := sql.Open("sqlite3", "file:transactions.db?cache=shared&mode=rwc")
	c.db = db
	if err != nil {
		return err
	}
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS transactions (
        id TEXT PRIMARY KEY,
        date TEXT,
        description TEXT,
        amount INT,
		category TEXT
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	data, err := loadDataFromServer[account.Data]("http://localhost:8080/data", c)
	if err != nil {
		return err
	}
	err = SaveDataToDB(c.db, data)
	return err
}

func SaveDataToDB(db *sql.DB, data account.Data) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM transactions")
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO transactions(id, date, description, amount, category) VALUES (?,?,?,?,?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	for _, t := range data.Transactions {
		_, err := stmt.Exec(t.ID, t.Date, t.Description, t.Amount, t.Category)
		fmt.Println("adding", t)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (c *Client) QueryTransactionsAndUpdate(info TransactionFilterInfo) error {
	db, err := sql.Open("sqlite3", "file:transactions.db?cache=shared&mode=rwc")
	if err != nil {
		return err
	}
	var filters []interface{}
	var filters2 []string

	if info.ID != "" {
		filters2 = append(filters2, "id = ?")
		filters = append(filters, info.ID)
	}
	if info.Date != "" {
		filters2 = append(filters2, "date = ?")
		filters = append(filters, info.Date)
	}
	if info.Description != "" {
		filters2 = append(filters2, "description = ?")
		filters = append(filters, info.Description)
	}
	if info.Amount != 0 {
		filters2 = append(filters2, "amount = ?")
		filters = append(filters, info.Amount)
	}
	if info.Category != "" {
		filters2 = append(filters2, "category = ?")
		filters = append(filters, info.Category)
	}
	//TODO: make it so amount can be greater less than or equal to
	query := `SELECT id, date, description, amount, category FROM transactions `
	query += "WHERE " + strings.Join(filters2, " AND ")

	rows, err := db.Query(query, filters...)
	if err != nil {
		return err
	}
	var NewTransactionList []account.Transaction
	for rows.Next() {
		var id string
		var date string
		var description string
		var amount int64
		var category string

		rows.Scan(&id, &date, &description, &amount, &category)
		t := account.NewTransaction(date, description, amount, category)
		NewTransactionList = append(NewTransactionList, t)
	}
	fmt.Println("NEW TRANSACTION LIST")
	fmt.Println(NewTransactionList)
	c.a.SetTransactionData(NewTransactionList)
	return nil
}
