package ui

import (
	"client/account"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
)

type AccountInterface interface {
	GetData() account.Data
}

type UIApp struct {
	fyneWindow fyne.Window
	fyneApp    fyne.App
	accountI   AccountInterface
}

func NewUIApp(a AccountInterface) *UIApp {
	result := UIApp{}
	result.fyneApp = app.New()
	result.fyneWindow = result.fyneApp.NewWindow("test")
	result.accountI = a
	return &result
}

func (ui *UIApp) ResizeWindow(width int, height int) {
	ui.fyneWindow.Resize(fyne.NewSize(width, height))
}

func (ui *UIApp) StartApp() {
	ui.fyneWindow.ShowAndRun()
}

func (ui *UIApp) loadDataIntoUI() {
	//transactionList := ui.account.GetData

	//list := widget.NewList(ui.account.getData().Transactions, createItem, updateItem)
}
