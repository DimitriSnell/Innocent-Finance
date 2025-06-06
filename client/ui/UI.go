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
	SetOnSyncStatusChanged(cb func(synced bool))
	GetDonatorNameByID(id string) string
}

type UIApp struct {
	fyneWindow fyne.Window
	fyneApp    fyne.App
	accountI   AccountInterface
	Width      float32
	tabMap     map[string]*tab
	tabList    []*tab
	currentTab string
	isSynced   bool
	tabs       *container.AppTabs
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

func (ui *UIApp) GetSynced() bool {
	return ui.isSynced
}

func (ui *UIApp) SetSynced(b bool) {
	ui.isSynced = b
	ui.LoadDataIntoUI()
}

func NewUIApp(a AccountInterface) *UIApp {
	result := &UIApp{}
	result.fyneApp = app.New()
	result.fyneWindow = result.fyneApp.NewWindow("test")
	result.accountI = a
	baseStruct := NewTab(DB.TransactionFilterInfo{}, 0, "base tab")
	result.tabMap = make(map[string]*tab)
	result.tabMap["base tab"] = baseStruct
	result.currentTab = "base tab"
	result.tabList = append(result.tabList, baseStruct)
	var amount int64
	amount = 100
	info := DB.TransactionFilterInfo{
		ID:          "",
		Date:        "",
		Description: "",
		Amount:      &amount,
		Category:    "",
	}
	baseStruct2 := NewTab(info, 0, "base tab2")
	result.tabMap["base tab2"] = baseStruct2
	result.tabList = append(result.tabList, baseStruct2)
	//set callback function
	a.SetOnSyncStatusChanged(func(b bool) {
		fmt.Println("TEST FROM CALLBACk")
		result.SetSynced(b)
	})
	result.fyneWindow.Resize(fyne.NewSize(1200, 600))
	return result
}

func (ui *UIApp) ResizeWindow(width float32, height float32) {
	ui.fyneWindow.Resize(fyne.NewSize(width, height))
	ui.Width = width
}

func (ui *UIApp) StartApp() {
	ui.fyneWindow.ShowAndRun()
}

func (ui *UIApp) createNewTabPopup() {
	id := widget.NewEntry()
	date := widget.NewEntry()
	date.SetPlaceHolder("ex 2024-10-05")
	date.Resize(fyne.NewSize(100, 40))
	date2 := widget.NewEntry()
	date2.SetPlaceHolder("ex 2024-10-05")
	date2.Resize(fyne.NewSize(100, 40))
	date2.Hide()
	DateOoperatorSelect := widget.NewSelect([]string{"Like", "<", "=", ">", "Between"}, func(selected string) {
		if selected == "Between" {
			date2.Show()
		} else {
			date2.Hide()
		}

	})
	DateOoperatorSelect.Resize(fyne.NewSize(500, 40))
	description := widget.NewEntry()
	amountEntry := widget.NewEntry()
	operatorSelect := widget.NewSelect([]string{"<", "=", ">"}, func(selected string) {
		fmt.Println("Operator selected:", selected)
	})
	operatorSelect.SetSelected("=") // default operator

	// Put amountEntry and operatorSelect side by side
	amountContainer := container.NewHBox(amountEntry, operatorSelect)
	DateContainer := container.NewGridWithColumns(3, container.NewStack(date), container.NewStack(DateOoperatorSelect), container.NewStack(date2))
	category := widget.NewEntry()
	var dialog *widget.PopUp

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ID", Widget: id},
			{Text: "Date", Widget: DateContainer},
			{Text: "Description", Widget: description},
			{Text: "Amount", Widget: amountContainer},
			{Text: "Category", Widget: category},
		},

		OnSubmit: func() {
			var amount *int64
			var temp int64
			fmt.Sscanf(amountEntry.Text, "%d", &temp)
			// Create TransactionFilterInfo
			amount = &temp
			if amountEntry.Text == "" {
				amount = nil
			}
			info := DB.TransactionFilterInfo{
				ID:          id.Text,
				Date:        date.Text,
				Description: description.Text,
				Amount:      amount,
				Category:    category.Text,
				Op:          operatorSelect.Selected,
				SecondDate:  date2.Text,
				DateOp:      DateOoperatorSelect.Selected,
			}
			if DateOoperatorSelect.Selected == "Between" && (date2.Text == "" || date.Text == "") || (DateOoperatorSelect.Selected == "" && date.Text != "") {
				fmt.Println("Please pick a second date")
			} else {
				// Create new tab
				newTitle := fmt.Sprintf("Tab %d", len(ui.tabList)+1)
				newTab := NewTab(info, len(ui.tabList), newTitle)
				ui.tabMap[newTitle] = newTab
				ui.tabList = append(ui.tabList, newTab)
				ui.currentTab = newTitle
				// Reload UI to rebuild tabs
				dialog.Hide()
				ui.LoadDataIntoUI()
			}

		},
		OnCancel: func() {
			dialog.Hide()
		},
		SubmitText: "Create",
		CancelText: "Cancel",
	}
	dialog = widget.NewModalPopUp(container.NewVBox(form), ui.fyneWindow.Canvas())
	dialog.Resize(fyne.NewSize(400, 300))
	dialog.Show()
}

func validateFilterForm() bool {
	return false
}

func (ui *UIApp) RefreshTabContent() {
	//set new tab
	var syncText string
	if ui.GetSynced() {
		syncText = "Synced with server"
	} else {
		syncText = "Not synced"
	}

	syncStatus := widget.NewLabel(syncText)
	headerBar := container.NewVBox(syncStatus)
	fmt.Println(ui.currentTab)
	header, list, err := ui.tabMap[ui.currentTab].CreateAndReturnUIContext(ui.accountI)
	if err != nil {
		fmt.Println("ERROR CREATING UI CONTEXT")
		return
	}
	content := container.NewVScroll(list)
	content.SetMinSize(fyne.NewSize(1200, 600))
	fixedHeightContainer := container.NewVBox(ui.tabs, header, content, headerBar)
	minWidthRect := canvas.NewRectangle(color.Transparent)
	minWidthRect.SetMinSize(fyne.NewSize(350, 10)) // 300px wide, 10px tall
	leftPanel := container.NewVBox(
		minWidthRect,
		widget.NewLabel("Left Panel"),
		layout.NewSpacer(), // This makes the left panel expand to fill available space
	)
	//set scroll offset to bottom then check if theres a scroll position saved
	totalHeight := float32(list.Length()) * list.MinSize().Height
	list.ScrollToOffset(totalHeight)
	if ui.tabMap[ui.currentTab].GetOffset() != -1 {
		list.ScrollToOffset(ui.tabMap[ui.currentTab].GetOffset())
	}

	// Create split container
	split := container.NewHSplit(leftPanel, fixedHeightContainer)
	split.SetOffset(0.2)
	split.Refresh()
	ui.fyneWindow.SetContent(split)
	ui.fyneWindow.Resize(fyne.NewSize(1200, 600)) // Make sure window is big enough
}

func (ui *UIApp) LoadDataIntoUI() error {
	var syncText string
	if ui.GetSynced() {
		syncText = "Synced with server"
	} else {
		syncText = "Not synced"
	}

	syncStatus := widget.NewLabel(syncText)
	headerBar := container.NewVBox(syncStatus)

	var tabBarItems []*container.TabItem
	for _, t := range ui.tabList {
		tabBarItems = append(tabBarItems, container.NewTabItem(t.title, container.NewWithoutLayout()))
	}
	// Add a final tab with a "+" button for adding a new tab
	addTabButton := widget.NewButton("+", func() {
		// Optionally add logic to create a new tab dynamically here
	})
	addTabButtonTab := container.NewTabItem("+", container.NewCenter(addTabButton))
	tabBarItems = append(tabBarItems, addTabButtonTab)

	tabs := container.NewAppTabs(tabBarItems...)
	header, list, err := ui.tabMap[ui.currentTab].CreateAndReturnUIContext(ui.accountI)
	tabs.SetTabLocation(container.TabLocationTop)
	ui.tabs = tabs
	//sets selected tab to current tab needed for when creating a new tab
	for i, t := range ui.tabList {
		if t.title == ui.currentTab {
			tabs.SelectIndex(i)
			break
		}
	}

	if err != nil {
		fmt.Println(err)
		return err
	}
	tabs.OnSelected = func(tab *container.TabItem) {
		tabString := string(tab.Text)
		if tabString == "+" {
			ui.createNewTabPopup()
			for i, t := range ui.tabList {
				if t.title == ui.currentTab {
					tabs.SelectIndex(i)
					break
				}
			}
			return
		}
		ui.tabMap[ui.currentTab].SetOffset(list.GetScrollOffset())
		ui.currentTab = tabString
		ui.RefreshTabContent()
	}

	//vscroll first is necessary for some reason
	content := container.NewVScroll(list)
	content.SetMinSize(fyne.NewSize(1200, 600))
	fixedHeightContainer := container.NewVBox(tabs, header, content, headerBar)
	fmt.Println(fixedHeightContainer)
	minWidthRect := canvas.NewRectangle(color.Transparent)
	minWidthRect.SetMinSize(fyne.NewSize(350, 200)) // 300px wide, 10px tall
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
	ui.fyneWindow.SetContent(split)
	ui.fyneWindow.Resize(fyne.NewSize(1200, 600)) // Make sure window is big enough
	return nil
}
