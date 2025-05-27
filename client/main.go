package main

import (
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
	client.GetAccount().AppendTransaction(account.NewTransaction("2023-10-01", "Test Transaction", 100, "Test Category"))
	fmt.Println(client.CheckServerSync())
	client.SyncServer()
	//UI := ui.NewUIApp(a)
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
