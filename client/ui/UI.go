package UI

import (
	"client/account"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type AccountInterface interface {
	GetData() account.Data
}

type UIApp struct {
	fyneWindow fyne.Window
	fyneApp    fyne.App
	accountI   AccountInterface
	Width      float32
}

type RightClickLabel struct {
	widget.BaseWidget
	content      fyne.CanvasObject
	onRightClick func()
}

func NewRightClickLabel(obj fyne.CanvasObject, onRightClick func()) *RightClickLabel {
	result := RightClickLabel{}
	result.onRightClick = onRightClick
	result.content = obj
	result.ExtendBaseWidget(&result)
	return &result
}

func (r *RightClickLabel) TappedSecondary(ev *fyne.PointEvent) {
	if r.onRightClick != nil {
		r.onRightClick()
	}
}

func (r *RightClickLabel) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(r.content)
}

func NewUIApp(a AccountInterface) *UIApp {
	result := UIApp{}
	result.fyneApp = app.New()
	result.fyneWindow = result.fyneApp.NewWindow("test")
	result.accountI = a
	return &result
}

func (ui *UIApp) ResizeWindow(width float32, height float32) {
	ui.fyneWindow.Resize(fyne.NewSize(width, height))
	ui.Width = width
}

func (ui *UIApp) StartApp() {
	ui.fyneWindow.ShowAndRun()
}

func (ui *UIApp) LoadDataIntoUI() {
	transactionList := ui.accountI.GetData().Transactions
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
	fmt.Println(list.Size().Width)
	//vscroll first is necessary for some reason
	content := container.NewVScroll(list)
	content.SetMinSize(fyne.NewSize(200, 200))
	fixedHeightContainer := container.NewVBox(header, content)
	//rect := canvas.NewRectangle(color.Transparent)
	//wrapped := NewRightClickLabel(rect, func() {
	//	fmt.Println("Right click detected on window!")
	//})
	minWidthRect := canvas.NewRectangle(color.Transparent)
	minWidthRect.SetMinSize(fyne.NewSize(250, 10)) // 300px wide, 10px tall
	leftPanel := container.NewVBox(
		minWidthRect,
		widget.NewLabel("Left Panel"),
		layout.NewSpacer(), // This makes the left panel expand to fill available space
	)
	// Create split container

	split := container.NewHSplit(leftPanel, fixedHeightContainer)
	split.SetOffset(0.2)
	split.Refresh()
	//allin := container.NewStack(wrapped, split)
	ui.fyneWindow.Resize(fyne.NewSize(900, 600)) // Make sure window is big enough
	ui.fyneWindow.SetContent(split)
}
