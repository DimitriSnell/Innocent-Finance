package clientapp

import (
	"bytes"
	ui "client/UI"
	"client/account"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	a  *account.Account
	ui *ui.UIApp
}

func NewClient() (*Client, error) {
	result := Client{}
	data, err := loadDataFromServer[account.Data]("http://localhost:8080/data")
	if err != nil {
		return &result, err
	}
	a, err := account.NewAccount(data)
	if err != nil {
		return &result, err
	}
	result.a = a
	result.ui = ui.NewUIApp(a)
	return &result, nil
}

func (c *Client) GetUI() *ui.UIApp {
	return c.ui
}

func (c *Client) CheckServerSync() (bool, error) {
	resp, err := http.Get("http://localhost:8080/sync-token")
	if err != nil {
		fmt.Println("Error loading data from server when checking sync", err)
		return false, err
	}
	defer resp.Body.Close()
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
	updatedServerData, err := c.PullToSync(c.a.GetData())
	if err != nil {
		return err
	}
	updatedServerData, err = c.pushToSync(updatedServerData, c.a.GetChanges())
	if err != nil {
		return err
	}
	c.a.SetData(updatedServerData)

	WriteDataToFile(c.a.GetData(), "data/data.json")
	return nil
}

func (c *Client) pushToSync(data account.Data, changes account.Changes) (account.Data, error) {
	result := data
	for _, v := range changes.AddedTransactions {
		result.Transactions = append(result.Transactions, v)
	}
	jsonData, err := json.Marshal(changes)
	if err != nil {
		return account.Data{}, err
	}
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
	result.SyncToken = st
	fmt.Println(result.SyncToken)
	return result, nil
}

func (c *Client) PullToSync(data account.Data) (account.Data, error) {
	dataFromServer, err := loadDataFromServer[account.Changes]("http://localhost:8080/pull-changes")
	if err != nil {
		return account.Data{}, err
	}
	updatedSyncToken, err := loadDataFromServer[int64]("http://localhost:8080/sync-token")
	if err != nil {
		return account.Data{}, err
	}
	if updatedSyncToken <= data.SyncToken {
		return data, nil
	}
	result := data
	for _, v := range dataFromServer.AddedTransactions {
		result.Transactions = append(result.Transactions, v)
	}
	return result, nil
}

func loadDataFromServer[T any](url string) (T, error) {
	var result T
	resp, err := http.Get(url)

	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	fmt.Println("response status:", resp.Status)
	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("failed to load data: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}
