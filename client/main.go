package main

import (
	DB "client/DB"
	"client/account"
	clientapp "client/clientApp"
	"fmt"
)

var pathToFile = "./data/data.json"

func main() {
	fmt.Println("hello world")
	/*a, err := account.NewAccount()
	if err != nil {
		fmt.Println("Error creating account:", err)
		return
	}
	fmt.Println(a)
	err = account.WriteDataToFile(a.GetData(), pathToFile)
	if err != nil {
		fmt.Println("Error writing data to file:", err)
		return
	}*/
	client, err := clientapp.NewClient()
	if err != nil {
		fmt.Println("Error creating account:", err)
		if err.Error() == "unauthorized" {
			fmt.Println("unauthorized detected")
		}
		return
	}

	fmt.Println(client)
	fmt.Println(client.CheckServerSync())
	Tlist := []account.Transaction{}
	client.AddTransactions(Tlist)
	DeleteList := []account.Transaction{}
	toDelete, err := DB.QueryTransactionByUID("9a908847-7d77-4e67-ae63-c433d7449c9f")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("TO DELETE ID", toDelete.ID)
	DeleteList = append(DeleteList, toDelete)
	client.DeleteTransactions(DeleteList)
	fmt.Println(client.CheckServerSync())
	client.SyncServer()
	//UI := ui.NewUIApp(a)
	var amount int64
	amount = 100
	info := DB.TransactionFilterInfo{
		ID:          "",
		Date:        "",
		Description: "",
		Amount:      &amount,
		Category:    "",
	}
	err = client.QueryTransactionsAndUpdate(info)
	if err != nil {
		fmt.Println(err)
		return
	}
	client.GetUI().LoadDataIntoUI()
	client.GetUI().ResizeWindow(500, 500)
	client.GetUI().StartApp()
	err = clientapp.DestroyDataFile(pathToFile)
	if err != nil {
		fmt.Println("Error destroying data file:", err)
		return
	}

	//fmt.Println(a)
}
