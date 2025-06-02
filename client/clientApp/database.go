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

// THIS FUNCTION IS ONLY USED ON SERVER STARTUP AS IT DELETES ALL DATA IN DATABASE AND CREATES A NEW DATABASE FROM THE DATA ARGUMENT
// THIS DATA ARGUMENT SHOULD ALWAYS BE THE MOST UP TO DATE DATA FROM THE SERVER
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
	//TODO add other parts of data like balance and syncToken
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

func QueryTransactionByUID(UID string) (account.Transaction, error) {
	var result account.Transaction
	db, err := sql.Open("sqlite3", "file:transactions.db?cache=shared&mode=rwc")
	if err != nil {
		return result, err
	}
	defer db.Close()
	query := `SELECT id, date, description, amount, category FROM transactions WHERE id = ?`
	rows, err := db.Query(query, UID)
	if err != nil {
		return result, err
	}
	count := 0
	for rows.Next() {
		count++
		if count > 1 {
			return result, fmt.Errorf("more than one transaction found with the same uid")
		}
		var id string
		var date string
		var description string
		var amount int64
		var category string

		err := rows.Scan(&id, &date, &description, &amount, &category)
		if err != nil {
			return result, err
		}
		result = account.NewTransaction(date, description, amount, category)
		result.ID = id
	}
	return result, nil
}

func QueryTransaction(info TransactionFilterInfo) ([]account.Transaction, error) {
	var result []account.Transaction
	db, err := sql.Open("sqlite3", "file:transactions.db?cache=shared&mode=rwc")
	if err != nil {
		return result, err
	}
	defer db.Close()
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
		return result, err
	}
	for rows.Next() {
		var id string
		var date string
		var description string
		var amount int64
		var category string

		rows.Scan(&id, &date, &description, &amount, &category)
		t := account.NewTransaction(date, description, amount, category)
		result = append(result, t)
	}
	return result, nil
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
