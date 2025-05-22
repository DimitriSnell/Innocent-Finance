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
	http.HandleFunc("/data", handleData)
	//http.HandleFunc("/add", handleAddTransaction)
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
