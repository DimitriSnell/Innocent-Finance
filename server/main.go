package main

import (
	"fmt"
	"project/handlers"
)

var a *Account

func main() {
	fmt.Println("hello world")
	handlers.Test()
	d, err := loadData()
	if err != nil {
		fmt.Println("error loading data!")
	}
	a = newAccount(d)
	fmt.Println(a.data.Budget)
	serverInit()
}
