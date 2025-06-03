package UI

import (
	DB "client/DB"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type tab struct {
	num    int
	title  string
	Filter DB.TransactionFilterInfo
}

func NewTab(F DB.TransactionFilterInfo, n int, t string) *tab {
	result := &tab{}
	result.num = n
	result.title = t
	result.Filter = F
	return result
}

func (t *tab) CreateAndReturnUIContext() (*fyne.Container, *widget.List, error) {
	transactionList, err := DB.QueryTransaction(t.Filter)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("length")
	fmt.Println(len(transactionList))
	header := container.NewGridWithColumns(9,
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Category", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Ammount", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Date", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewLabelWithStyle("Description", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
	)

	list := widget.NewList(
		func() int {
			return len(transactionList)
		},
		func() fyne.CanvasObject {
			base := container.NewGridWithColumns(9,
				layout.NewSpacer(),
				canvas.NewText("", color.Black), // Category
				layout.NewSpacer(),
				canvas.NewText("", color.Black), // Ammount
				layout.NewSpacer(),
				canvas.NewText("", color.Black), // Date
				layout.NewSpacer(),
				canvas.NewText("", color.Black), // Description
				layout.NewSpacer(),
			)
			rightClickWrap := NewRightClickLabel(base, func() {
				fmt.Println("test right click")
			})
			return rightClickWrap
		},
		func(li widget.ListItemID, co fyne.CanvasObject) {
			//co.(*widget.Label).SetText(fmt.Sprintf("Category: %s, Date: %s, Ammount: %d, Description: %s", transactionList[li].Category, transactionList[li].Date, transactionList[li].Amount, transactionList[li].Description))
			rightCickable := co.(*RightClickLabel)
			items := rightCickable.content.(*fyne.Container).Objects
			//items := co.(*fyne.Container).Objects
			if transactionList[li].Amount < 0 {
				items[1].(*canvas.Text).Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
				items[3].(*canvas.Text).Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
				items[5].(*canvas.Text).Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
				items[7].(*canvas.Text).Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
			} else {
				items[1].(*canvas.Text).Color = color.RGBA{R: 42, G: 168, B: 65, A: 255}
				items[3].(*canvas.Text).Color = color.RGBA{R: 42, G: 168, B: 65, A: 255}
				items[5].(*canvas.Text).Color = color.RGBA{R: 42, G: 168, B: 65, A: 255}
				items[7].(*canvas.Text).Color = color.RGBA{R: 42, G: 168, B: 65, A: 255}
			}
			items[1].(*canvas.Text).Text = transactionList[li].Category
			items[3].(*canvas.Text).Text = fmt.Sprintf("%d", transactionList[li].Amount)
			items[5].(*canvas.Text).Text = transactionList[li].Date
			items[7].(*canvas.Text).Text = transactionList[li].Description
			// Refresh texts after setting the text
			for _, idx := range []int{1, 3, 5, 7} {
				canvas.Refresh(items[idx].(*canvas.Text))
			}
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		fmt.Println("test?")
	}
	return header, list, nil
}
