package clientapp

import (
	"bytes"
	ui "client/UI"
	"client/account"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

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
	//TODO: here is where data is initally stored in memory eventually make it so only partial amounts of data and stored and the rest
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
	err = result.initDatabase()
	if err != nil {
		return &result, err
	}
	return &result, nil
}

func (c *Client) AddTransactions(transactionList []account.Transaction) error {
	stmt, err := c.db.Prepare("INSERT INTO transactions(id, date, description, amount, category) VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, t := range transactionList {
		_, err := stmt.Exec(t.ID, t.Date, t.Description, t.Amount, t.Category)
		c.a.AppendTransaction(t)
		fmt.Println("adding", t)
		if err != nil {
			return err
		}
	}
	for _, t := range c.GetAccount().GetData().Transactions {
		fmt.Println(t)
	}
	fmt.Println("len1")
	fmt.Println(len(c.GetAccount().GetData().Transactions))
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
	// pushes changes to server NOTE: does not POST full data only POSTS changes
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
	fmt.Println(result.SyncToken)
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
	//adds added transactions from the server to memory and data
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
