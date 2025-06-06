package DB

import (
	"client/account"
	"database/sql"
	"fmt"
	"strings"
)

type TransactionFilterInfo struct {
	ID          string
	Date        string
	DateOp      string
	SecondDate  string
	Description string
	Amount      *int64
	Op          string
	Category    string
	DonatorID   string
}

func InitDatabase(data account.Data) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:transactions.db?cache=shared&mode=rwc")
	if err != nil {
		return nil, err
	}
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS transactions (
        id TEXT PRIMARY KEY,
        date TEXT,
        description TEXT,
        amount INT,
		category TEXT,
		donatorid TEXT
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}
	createTableSQL = `
    CREATE TABLE IF NOT EXISTS donators (
        id TEXT PRIMARY KEY,
        name TEXT
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}
	err = SaveDataToDB(db, data)
	return db, err
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
	_, err = tx.Exec("DELETE FROM donators")
	if err != nil {
		tx.Rollback()
		return err
	}
	//TODO add other parts of data like balance and syncToken
	stmt, err := tx.Prepare("INSERT INTO transactions(id, date, description, amount, category, donatorid) VALUES (?,?,?,?,?,?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	for _, t := range data.Transactions {
		_, err := stmt.Exec(t.ID, t.Date, t.Description, t.Amount, t.Category, t.DonatorID)
		fmt.Println("adding", t)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	stmt, err = tx.Prepare("INSERT INTO donators(id, name) VALUES (?,?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, d := range data.Donators {
		_, err := stmt.Exec(d.ID, d.Name)
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
	query := `SELECT id, date, description, amount, category, donatorid FROM transactions WHERE id = ?`
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
		var donatorid string

		err := rows.Scan(&id, &date, &description, &amount, &category, &donatorid)
		if err != nil {
			return result, err
		}

		if donatorid == "" {
			donatorid = "anonymous"
		}
		result = account.NewTransaction(date, description, amount, category, donatorid)
		result.ID = id
	}
	return result, nil
}

func QueryDonatorList() ([]account.Donator, error) {
	var result []account.Donator
	db, err := sql.Open("sqlite3", "file:transactions.db?cache=shared&mode=rwc")
	if err != nil {
		return result, err
	}
	query := `SELECT id, name FROM donators`
	rows, err := db.Query(query)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var id string
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			return result, err
		}
		d := account.NewDonator(name)
		d.ID = id
		result = append(result, d)
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
		if info.DateOp == "" || info.DateOp == "Like" {
			fmt.Println(info.Date, info.DateOp, info.SecondDate)
			filters2 = append(filters2, "date LIKE ?")
			filters = append(filters, "%"+info.Date+"%")
		} else if info.DateOp != "Between" {
			filters2 = append(filters2, "date "+info.DateOp+" ?")
			filters = append(filters, info.Date)
		} else {
			filters2 = append(filters2, "date >= ? AND date <= ?")
			filters = append(filters, info.Date, info.SecondDate)
		}
	}
	if info.Description != "" {
		filters2 = append(filters2, "description LIKE ?")
		filters = append(filters, "%"+info.Description+"%")
	}
	if info.Amount != nil {
		if info.Op == "" {
			filters2 = append(filters2, "amount = ?")
			filters = append(filters, info.Amount)
		} else {
			filters2 = append(filters2, "amount "+info.Op+" ?")
			filters = append(filters, info.Amount)
		}
	}
	if info.Category != "" {
		filters2 = append(filters2, "category LIKE ?")
		filters = append(filters, "%"+info.Category+"%")
	}
	if info.DonatorID != "" {
		filters2 = append(filters2, "donatorid = ?")
		filters = append(filters, info.DonatorID)
	}
	//TODO: make it so amount can be greater less than or equal to
	query := `SELECT id, date, description, amount, category, donatorid FROM transactions `
	if len(filters2) > 0 {
		query += "WHERE " + strings.Join(filters2, " AND ")
	}
	fmt.Println(query)
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
		var donatorid string

		rows.Scan(&id, &date, &description, &amount, &category, &donatorid)
		t := account.NewTransaction(date, description, amount, category, donatorid)
		result = append(result, t)
	}
	return result, nil
}
