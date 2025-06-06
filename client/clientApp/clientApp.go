package clientapp

import (
	"bytes"
	DB "client/DB"
	ui "client/UI"
	"client/account"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	a     *account.Account
	ui    *ui.UIApp
	token string
	db    *sql.DB
}

func NewClient() (*Client, error) {
	result := Client{}
	result.token = "1234"
	//TODO: here is where data is initally stored in memory eventually make it so only partial amounts of data and stored in memory (for lazy loading) and the rest
	//TODO: is only put in the local database
	data, err := loadDataFromServer[account.Data]("http://localhost:8080/data", &result)
	if err != nil {
		return &result, err
	}
	a, err := account.NewAccount(data)
	if err != nil {
		return &result, err
	}
	result.a = a
	result.ui = ui.NewUIApp(a)
	tempDb, err := DB.InitDatabase(data)
	if err != nil {
		return &result, err
	}
	result.db = tempDb
	return &result, nil
}

func (c *Client) ClearChanges() {
	c.a.ClearChanges()
}

func (c *Client) DeleteTransactions(transactionList []account.Transaction) error {
	stmt, err := c.db.Prepare("DELETE FROM transactions WHERE id = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()
	fmt.Println("IN DELETE TRANSACTION")
	for _, t := range transactionList {
		_, err := stmt.Exec(t.ID)
		if err != nil {
			return err
		}
		fmt.Println("IN DELETE TRANSACTION2")
		c.a.DeleteTransactionFromMemory(t)
	}
	//run callback to set sync status false
	c.a.MarkUnsynced()
	return nil
}

func (c *Client) AddTransactions(transactionList []account.Transaction) error {
	stmt, err := c.db.Prepare("INSERT INTO transactions(id, date, description, amount, category, donatorid) VALUES (?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, t := range transactionList {
		_, err := stmt.Exec(t.ID, t.Date, t.Description, t.Amount, t.Category, t.DonatorID)
		if err != nil {
			return err
		}
		c.a.AppendTransaction(t)
	}
	//run callback to set sync status false
	c.a.MarkUnsynced()
	return nil
}

func (c *Client) AddDonators(donatorList []account.Donator) error {
	stmt, err := c.db.Prepare("INSERT INTO donators(id, name) VALUES (?,?)")
	if err != nil {
		return err
	}
	for _, d := range donatorList {
		_, err := stmt.Exec(d.ID, d.Name)
		if err != nil {
			return err
		}
		c.a.AppendDonator(&d)
	}
	c.a.MarkUnsynced()
	return nil
}

func (c *Client) GetUI() *ui.UIApp {
	return c.ui
}

func (c *Client) GetAccount() *account.Account {
	return c.a
}

func (c *Client) CheckServerSync() (bool, error) {
	req, err := http.NewRequest("GET", "http://localhost:8080/sync-token", nil)
	if err != nil {
		fmt.Println("Error creating request", err)
		return false, err
	}
	req.Header.Set("Authorization-Token", c.token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error loading data from server when checking sync", err)
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return false, fmt.Errorf("unauthorized")
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error loading data from server when checking sync", resp.Status)
		return false, nil
	}
	var st int64
	if err := json.NewDecoder(resp.Body).Decode(&st); err != nil {
		fmt.Println("Error decoding sync token from server", err)
		return false, err
	}
	if st != c.a.GetData().SyncToken {
		fmt.Println("server is not in sync with client")
		return false, nil
	}
	return true, nil
}

func (c *Client) SyncServer() error {
	//pulls and saves any updated data from the server then adds and
	//pushes changes to server NOTE: does not POST full data only POSTS changes
	//TODO potentially remove Data from return of pull/pushToSync
	updatedServerData, err := c.PullToSync(c.a.GetData())
	if err != nil {
		return err
	}
	updatedServerData, err = c.pushToSync(updatedServerData, c.a.GetChanges())
	if err != nil {
		return err
	}
	//c.a.SetData(updatedServerData)

	WriteDataToFile(c.a.GetData(), "data/data.json")
	c.a.NotifySyncStatusChanged(true)
	return nil
}

func (c *Client) pushToSync(data account.Data, changes account.Changes) (account.Data, error) {
	result := data
	jsonData, err := json.Marshal(changes)
	if err != nil {
		return account.Data{}, err
	}
	//TODO add authorization to POST request
	resp, err := http.Post("http://localhost:8080/push", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("error pushing data to server")
		return account.Data{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return account.Data{}, fmt.Errorf("failed to post data: %s", resp.Status)
	}
	//sets syncToken of result to the servers updated syncToken
	var st int64
	if err := json.NewDecoder(resp.Body).Decode(&st); err != nil {
		fmt.Println("Error decoding sync token from server", err)
		return account.Data{}, err
	}
	c.a.SetSyncToken(st)
	c.ClearChanges()
	return result, nil
}

func (c *Client) PullToSync(data account.Data) (account.Data, error) {
	dataFromServer, err := loadDataFromServer[account.Changes]("http://localhost:8080/pull-changes", c)
	if err != nil {
		return account.Data{}, err
	}
	updatedSyncToken, err := loadDataFromServer[int64]("http://localhost:8080/sync-token", c)
	if err != nil {
		return account.Data{}, err
	}
	if updatedSyncToken <= data.SyncToken {
		return data, nil
	}
	result := data
	//adds added changed transactions from the server to memory and data
	c.AddTransactions(dataFromServer.AddedTransactions)
	return result, nil
}

func loadDataFromServer[T any](url string, c *Client) (T, error) {
	var result T
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request", err)
		return result, err
	}
	req.Header.Set("Authorization-Token", c.token)

	client := &http.Client{}

	resp, err := client.Do(req)

	//resp, err := http.Get(url)

	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	fmt.Println("response status:", resp.Status)

	if resp.StatusCode == http.StatusUnauthorized {
		return result, fmt.Errorf("unauthorized")
	}

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("failed to load data: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}

func (c *Client) QueryTransactionsAndUpdate(info DB.TransactionFilterInfo) error {
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
	if info.Amount != nil {
		filters2 = append(filters2, "amount = ?")
		filters = append(filters, info.Amount)
	}
	if info.Category != "" {
		filters2 = append(filters2, "category = ?")
		filters = append(filters, info.Category)
	}
	if info.DonatorID != "" {
		filters2 = append(filters2, "donatorid = ?")
		filters = append(filters, info.DonatorID)
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
		var donatorid string

		rows.Scan(&id, &date, &description, &amount, &category, &donatorid)
		t := account.NewTransaction(date, description, amount, category, donatorid)
		NewTransactionList = append(NewTransactionList, t)
	}
	fmt.Println("NEW TRANSACTION LIST")
	fmt.Println(NewTransactionList)
	c.a.SetTransactionData(NewTransactionList)
	return nil
}
