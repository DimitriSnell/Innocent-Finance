package main

import (
	"encoding/json"
	"fmt"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

var pathToFile = "./data/data.json"

func main() {
	fmt.Println("hello world")
	a, err := newAccount()
	if err != nil {
		fmt.Println("Error creating account:", err)
		return
	}
	fmt.Println(a)
	err = writeDataToFile(a.data, pathToFile)
	if err != nil {
		fmt.Println("Error writing data to file:", err)
		return
	}
	fmt.Println(a.checkServerSync())
	a.appendTransaction(newTransaction("2023-10-01", "Test Transaction", 100, "Test Category"))
	fmt.Println(a.checkServerSync())
	a.syncServer()
	app := app.New()
	w := app.NewWindow("Hello")
	w.SetContent(widget.NewLabel("Hello Fyne!"))
	w.Resize(fyne.NewSize(800, 500))
	w.ShowAndRun()

	err = destroyDataFile(pathToFile)
	if err != nil {
		fmt.Println("Error destroying data file:", err)
		return
	}

	fmt.Println(a)
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

func destroyDataFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}
