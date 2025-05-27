package main

import (
	ui "client/UI"
	"client/account"
	"fmt"
)

var pathToFile = "./data/data.json"

func main() {
	fmt.Println("hello world")
	a, err := account.NewAccount()
	if err != nil {
		fmt.Println("Error creating account:", err)
		return
	}
	fmt.Println(a)
	err = account.WriteDataToFile(a.GetData(), pathToFile)
	if err != nil {
		fmt.Println("Error writing data to file:", err)
		return
	}
	fmt.Println(a.CheckServerSync())
	//a.appendTransaction(newTransaction("2023-10-01", "Test Transaction", 100, "Test Category"))
	fmt.Println(a.CheckServerSync())
	a.SyncServer()
	UI := ui.NewUIApp(a)
	UI.LoadDataIntoUI()
	UI.ResizeWindow(500, 500)
	UI.StartApp()
	err = account.DestroyDataFile(pathToFile)
	if err != nil {
		fmt.Println("Error destroying data file:", err)
		return
	}

	fmt.Println(a)
}
