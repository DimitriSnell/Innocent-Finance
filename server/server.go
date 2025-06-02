package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

var (
	dataFile = "data/data.json"
	mu       sync.Mutex
)

func serverInit() {
	dataHandler := http.HandlerFunc(handleData)
	http.Handle("/data", authenticationMiddleware(dataHandler))
	http.HandleFunc("/sync-token", syncTokenHandler)
	//http.HandleFunc("/add", handleAddTransaction)
	http.HandleFunc("/push", handleSync)
	http.HandleFunc("/pull-changes", handlePull)
	fmt.Println("listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func loadData() (Data, error) {
	var d Data
	var err error
	bytes, err := os.ReadFile(dataFile)

	if err == nil {
		err = json.Unmarshal(bytes, &d)
	}

	return d, err
}

func handleData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a.data)
}

// returns sync token to client
func syncTokenHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("sync token handler")
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a.data.SyncToken)
}

func handleSync(w http.ResponseWriter, r *http.Request) {
	fmt.Println("sync handler")
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	var changes Changes
	err := json.NewDecoder(r.Body).Decode(&changes)
	if err != nil {
		http.Error(w, "Error decoding json", http.StatusBadRequest)
		return
	}
	//add changes from client and increment then return syncToken
	for _, v := range changes.AddedTransactions {
		a.data.Transactions = append(a.data.Transactions, v)
		a.data.SyncToken++
	}
	for _, v := range changes.DeletedTransactions {
		idx := FindIndexByID(v.ID, a.data.Transactions)
		fmt.Println(idx)
		if idx == -1 {
			break
		}
		newList := append(a.data.Transactions[:idx], a.data.Transactions[idx+1:]...)
		a.data.Transactions = newList
		a.data.SyncToken++
	}
	writeDataToFile(a.data, dataFile)
	json.NewEncoder(w).Encode(a.data.SyncToken)
}

func writeDataToFile(data Data, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func handlePull(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a.changes)

}
