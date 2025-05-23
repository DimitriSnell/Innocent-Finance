package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Budget struct {
	MonthlyLimit int64 `json:"monthly_limit"`
	TotalBudget  int64 `json:"total_budget"`
}

type Data struct {
	Transactions []Transaction `json:"transactions"`
	Budget       Budget        `json:"budget"`
	SyncToken    int64         `json:"sync_token"`
}

type Account struct {
	data    Data
	changes Changes
}

type Changes struct {
	AddedTransactions    []Transaction `json:"addedTransactions,omitempty"`
	DeletedTransactions  []Transaction `json:"deletedTransactions,omitempty"`
	ReplacedTransactions []Transaction `json:"replacedTransactions,omitempty"`
}

func newAccount() (*Account, error) {
	a := Account{}
	d, err := loadDataFromServer[Data]("http://localhost:8080/data")
	if err != nil {
		return nil, err
	}
	a.data = d
	a.changes = Changes{}
	return &a, nil
}

func (a *Account) appendTransaction(t Transaction) {
	a.data.Transactions = append(a.data.Transactions, t)
	a.changes.AddedTransactions = append(a.changes.AddedTransactions, t)
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

func (a *Account) checkServerSync() (bool, error) {
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
	if st != a.data.SyncToken {
		fmt.Println("server is not in sync with client")
		return false, nil
	}
	return true, nil
}

func (a *Account) syncServer() error {
	//pulls and saves any updated data from the server then adds and
	// pushes changes to server NOTE: does not POST full data only POSTS changes
	updatedServerData, err := PullToSync(a.data)
	updatedServerData, err = pushToSync(updatedServerData, a.changes)
	if err != nil {
		return err
	}
	a.data = updatedServerData
	writeDataToFile(a.data, "data/data.json")
	return nil
}

func pushToSync(data Data, changes Changes) (Data, error) {
	result := data
	for _, v := range changes.AddedTransactions {
		result.Transactions = append(result.Transactions, v)
	}
	jsonData, err := json.Marshal(changes)
	resp, err := http.Post("http://localhost:8080/push", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("error pushing data to server")
		return Data{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Data{}, fmt.Errorf("failed to post data: %s", resp.Status)
	}
	//sets syncToken of result to the servers updated syncToken
	var st int64
	if err := json.NewDecoder(resp.Body).Decode(&st); err != nil {
		fmt.Println("Error decoding sync token from server", err)
		return Data{}, err
	}
	result.SyncToken = st
	fmt.Println(result.SyncToken)
	return result, nil
}

func PullToSync(data Data) (Data, error) {
	dataFromServer, err := loadDataFromServer[Changes]("http://localhost:8080/pull-changes")
	if err != nil {
		return Data{}, err
	}
	updatedSyncToken, err := loadDataFromServer[int64]("http://localhost:8080/sync-token")
	if err != nil {
		return Data{}, err
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
