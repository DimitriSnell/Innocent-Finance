package UI

import (
	DB "client/DB"
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
	tabMap     map[string]*tab
	tabList    []*tab
	currentTab string
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
	baseStruct := NewTab(DB.TransactionFilterInfo{}, 0, "base tab")
	result.tabMap = make(map[string]*tab)
	result.tabMap["base tab"] = baseStruct
	result.currentTab = "base tab"
	result.tabList = append(result.tabList, baseStruct)
	info := DB.TransactionFilterInfo{
		ID:          "",
		Date:        "",
		Description: "",
		Amount:      100,
		Category:    "",
	}
	baseStruct2 := NewTab(info, 0, "base tab2")
	result.tabMap["base tab2"] = baseStruct2
	result.tabList = append(result.tabList, baseStruct2)
	return &result
}

func (ui *UIApp) ResizeWindow(width float32, height float32) {
	ui.fyneWindow.Resize(fyne.NewSize(width, height))
	ui.Width = width
}

func (ui *UIApp) StartApp() {
	ui.fyneWindow.ShowAndRun()
}

func (ui *UIApp) LoadDataIntoUI() error {

	/*tabs := container.NewHBox(
		widget.NewButton("All", func() {
			fmt.Println("All clicked")
			// You can update the list content here
		}),
		widget.NewButton("Expenses", func() {
			fmt.Println("Expenses clicked")
		}),
		widget.NewButton("Income", func() {
			fmt.Println("Income clicked")
		}),
	)*/

	var tabBarItems []*container.TabItem
	for _, t := range ui.tabList {
		tabBarItems = append(tabBarItems, container.NewTabItem(t.title, widget.NewLabel("test")))
	}
	tabs := container.NewAppTabs(tabBarItems...)
	header, list, err := ui.tabMap[ui.currentTab].CreateAndReturnUIContext()
	tabs.SetTabLocation(container.TabLocationTop)
	if err != nil {
		fmt.Println(err)
		return err
	}
	tabs.OnSelected = func(tab *container.TabItem) {
		tabString := string(tab.Text)
		ui.currentTab = tabString
		fmt.Println(ui.currentTab)
		header, list, err = ui.tabMap[ui.currentTab].CreateAndReturnUIContext()
		content := container.NewVScroll(list)
		content.SetMinSize(fyne.NewSize(200, 200))
		fixedHeightContainer := container.NewVBox(tabs, header, content)
		fmt.Println(fixedHeightContainer)
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
		//ui.fyneWindow.Resize(fyne.NewSize(900, 600)) // Make sure window is big enough
		ui.fyneWindow.SetContent(split)
	}
	//vscroll first is necessary for some reason
	content := container.NewVScroll(list)
	content.SetMinSize(fyne.NewSize(200, 200))
	fixedHeightContainer := container.NewVBox(tabs, header, content)
	fmt.Println(fixedHeightContainer)
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
	return nil
}
