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

type ContentType int

const (
	TransactionType ContentType = iota
	DonatorType
)

type tab struct {
	num    int
	title  string
	Filter DB.TransactionFilterInfo
	cType  ContentType
	offset float32
}

func NewTab(F DB.TransactionFilterInfo, n int, t string, ct ContentType) *tab {
	result := &tab{}
	result.num = n
	result.title = t
	result.Filter = F
	result.offset = -1
	result.cType = ct
	return result
}

func (t *tab) SetOffset(o float32) {
	t.offset = o
}

func (t *tab) GetOffset() float32 {
	return t.offset
}

func (t *tab) SetType(ct ContentType) {
	t.cType = ct
}

func (t *tab) GetType() ContentType {
	return t.cType
}

func (t *tab) CreateAndReturnUIContext(ai AccountInterface) (*fyne.Container, *widget.List, error) {
	var cont *fyne.Container
	var list *widget.List
	var err error
	switch t.cType {
	case TransactionType:
		cont, list, err = t.CreateTransactionContent(ai)
	case DonatorType:
		cont, list, err = t.CreateDonatorContent(ai)
	}
	return cont, list, err
}

func (t *tab) CreateDonatorContent(ai AccountInterface) (*fyne.Container, *widget.List, error) {
	DonatorList, err := DB.QueryDonatorList()
	if err != nil {
		return nil, nil, err
	}
	header := container.New(layout.NewAdaptiveGridLayout(3),
		widget.NewLabel(""),
		widget.NewLabelWithStyle("Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Amount donated this month", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)
	list := widget.NewList(
		func() int {
			return len(DonatorList)
		},
		func() fyne.CanvasObject {

			base := container.New(layout.NewGridLayoutWithColumns(3),
				canvas.NewText("", color.White), // empty first column
				canvas.NewText("", color.White), // Name
				canvas.NewText("", color.White), // Amount
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

			items[1].(*canvas.Text).Text = DonatorList[li].Name
			if DonatorList[li].Name == "" {
				items[1].(*canvas.Text).Text = "Unknown"
			}
			items[2].(*canvas.Text).Text = "Not Available"
			// Refresh texts after setting the text
		},
	)
	return header, list, nil
}

func (t *tab) CreateTransactionContent(ai AccountInterface) (*fyne.Container, *widget.List, error) {
	transactionList, err := DB.QueryTransaction(t.Filter)
	if err != nil {
		return nil, nil, err
	}
	header := container.New(layout.NewGridLayoutWithColumns(6),
		widget.NewLabel(""),
		widget.NewLabelWithStyle("Category", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Amount", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Date", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Description", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Donator", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	)

	list := widget.NewList(
		func() int {
			return len(transactionList)
		},
		func() fyne.CanvasObject {

			base := container.New(layout.NewGridLayoutWithColumns(6),
				canvas.NewText("", color.Black), // empty first column
				canvas.NewText("", color.Black), // Category
				canvas.NewText("", color.Black), // Amount
				canvas.NewText("", color.Black), // Date
				canvas.NewText("", color.Black), // Description
				canvas.NewText("", color.Black), // Donator
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
			/*if transactionList[li].Amount < 0 {
				items[1].(*canvas.Text).Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
				items[2].(*canvas.Text).Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
				items[3].(*canvas.Text).Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
				items[4].(*canvas.Text).Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
				items[5].(*canvas.Text).Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
			} else {
				items[1].(*canvas.Text).Color = color.RGBA{R: 42, G: 168, B: 65, A: 255}
				items[2].(*canvas.Text).Color = color.RGBA{R: 42, G: 168, B: 65, A: 255}
				items[3].(*canvas.Text).Color = color.RGBA{R: 42, G: 168, B: 65, A: 255}
				items[4].(*canvas.Text).Color = color.RGBA{R: 42, G: 168, B: 65, A: 255}
				items[5].(*canvas.Text).Color = color.RGBA{R: 42, G: 168, B: 65, A: 255}
			}
			items[1].(*canvas.Text).Text = transactionList[li].Category
			items[2].(*canvas.Text).Text = fmt.Sprintf("%d", transactionList[li].Amount)
			items[3].(*canvas.Text).Text = transactionList[li].Date
			items[4].(*canvas.Text).Text = transactionList[li].Description
			items[5].(*canvas.Text).Text = ai.GetDonatorNameByID(transactionList[li].DonatorID)*/

			for _, obj := range items {
				if t, ok := obj.(*canvas.Text); ok {
					t.Color = color.Black // default color
				}
			}

			if transactionList[li].Amount < 0 {
				for _, idx := range []int{1, 2, 3, 4, 5} {
					if t, ok := items[idx].(*canvas.Text); ok {
						t.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // red
					}
				}
			} else {
				for _, idx := range []int{1, 2, 3, 4, 5} {
					if t, ok := items[idx].(*canvas.Text); ok {
						t.Color = color.RGBA{R: 42, G: 168, B: 65, A: 255} // green
					}
				}
			}

			items[1].(*canvas.Text).Text = transactionList[li].Category
			items[2].(*canvas.Text).Text = fmt.Sprintf("%d", transactionList[li].Amount)
			items[3].(*canvas.Text).Text = transactionList[li].Date
			items[4].(*canvas.Text).Text = transactionList[li].Description
			items[5].(*canvas.Text).Text = ai.GetDonatorNameByID(transactionList[li].DonatorID)
			// Refresh texts after setting the text
			for _, idx := range []int{1, 2, 3, 4, 5} {
				items[idx].Refresh()
			}
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		fmt.Println("test?")
	}
	return header, list, nil
}
