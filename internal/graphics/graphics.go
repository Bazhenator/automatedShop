package graphics

import (
	"automatedShop/internal/services"
	"automatedShop/internal/services/dto"
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jung-kurt/gofpdf"
	"os"
	"path/filepath"
	"strconv"
)

const (
	successfulLoginMsg = "user logged in successfully!"
)

type AppManager struct {
	AuthService services.IAuthService
	ShopService services.IShopService
	UserLabel   *widget.Entry
}

func NewAppManager(s *services.Service) *AppManager {
	userLabel := widget.NewEntry()

	return &AppManager{
		AuthService: s.AuthService,
		ShopService: s.ShopService,
		UserLabel:   userLabel,
	}
}

func (m *AppManager) Run() {
	application := app.New()
	mainWindow := application.NewWindow("Shop Management System v.0.0.0")
	mainWindow.Resize(fyne.NewSize(300, 300))

	m.ShowLoginScreen(mainWindow)

	mainWindow.ShowAndRun()
}

// ShowLoginScreen shows login screen window to user
func (m *AppManager) ShowLoginScreen(window fyne.Window) {
	m.UserLabel.SetPlaceHolder("Username")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	errorLabel := widget.NewLabel("")

	loginButton := widget.NewButton("Login", func() {
		username := m.UserLabel.Text
		password := passwordEntry.Text

		if m.AuthService.AuthoriseUser(context.Background(), username, password) {
			dialog.ShowInformation("Authorized", successfulLoginMsg, window)
			m.ShowMainScreen(window, username)
		} else {
			errorLabel.SetText("Invalid username or password")
		}
	})

	registerButton := widget.NewButton("Register", func() {
		m.ShowRegisterScreen(window)
	})

	window.SetContent(container.NewStack(
		container.NewVBox(
			widget.NewLabel("Authorization"),
			m.UserLabel,
			passwordEntry,
			loginButton,
			registerButton,
			errorLabel,
		)))
}

// ShowRegisterScreen shows register screen window to users
func (m *AppManager) ShowRegisterScreen(window fyne.Window) {
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Username")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	errorLabel := widget.NewLabel("")

	registerButton := widget.NewButton("Register", func() {
		username := usernameEntry.Text
		password := passwordEntry.Text

		err := m.AuthService.RegisterUser(context.Background(), username, password)
		if err != nil {
			errorLabel.SetText("Registration failed: " + err.Error())
		} else {
			m.ShowLoginScreen(window)
		}
	})

	exitButton := widget.NewButton("Back", func() {
		m.ShowLoginScreen(window)
	})

	form := container.NewVBox(
		widget.NewLabel("Register"),
		usernameEntry,
		passwordEntry,
		registerButton,
		exitButton,
		errorLabel,
	)

	window.SetContent(form)
}

// ShowMainScreen shows main application's screen to user
func (m *AppManager) ShowMainScreen(window fyne.Window, loginEntry string) {
	loginLabel := widget.NewLabelWithStyle("user: "+loginEntry, fyne.TextAlignTrailing, fyne.TextStyle{Monospace: true})

	tabs := container.NewAppTabs(
		container.NewTabItem("Handbooks", m.ShowHandbooksScreen(window)),
		container.NewTabItem("Journals", m.ShowJournalsScreen(window)),
		container.NewTabItem("Reports", m.ShowReportsScreen(window)),
	)

	exitButton := widget.NewButton("Logout", func() {
		m.ShowLoginScreen(window)
	})

	mainLayout := container.NewBorder(loginLabel, exitButton, nil, nil, tabs)
	window.SetContent(mainLayout)
}

// ShowHandbooksScreen shows screen with handbooks' features to user
func (m *AppManager) ShowHandbooksScreen(window fyne.Window) fyne.CanvasObject {
	warehousesButton := widget.NewButton("Warehouses", func() {
		m.ShowWarehousesTable(window)
	})

	expenseItemsButton := widget.NewButton("Expense Items", func() {
		m.ShowExpenseItemsTable(window)
	})

	return container.NewVBox(
		widget.NewLabelWithStyle("Please, choose Handbook:", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}),
		warehousesButton,
		expenseItemsButton,
	)
}

// ShowWarehousesTable outputs data from warehouses table
func (m *AppManager) ShowWarehousesTable(window fyne.Window) {
	data, err := m.ShopService.ShowWarehousesTable(context.Background())
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	headers := []string{"id", "name", "quantity", "amount"}

	table := widget.NewTable(
		func() (int, int) { return len(data) + 1, len(headers) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			if id.Row == 0 {
				label.SetText(headers[id.Col])
				label.TextStyle = fyne.TextStyle{Monospace: true}
				return
			}

			row := id.Row - 1
			switch id.Col {
			case 0:
				label.SetText(strconv.Itoa(data[row].Id))
			case 1:
				label.SetText(data[row].Name)
			case 2:
				label.SetText(strconv.Itoa(data[row].Quantity))
			case 3:
				label.SetText(strconv.Itoa(data[row].Amount))
			}
			label.TextStyle = fyne.TextStyle{Monospace: true}
		},
	)

	table.SetColumnWidth(0, 50)  // ID
	table.SetColumnWidth(1, 200) // Name
	table.SetColumnWidth(2, 100) // Quantity
	table.SetColumnWidth(3, 100) // Amount

	tableContainer := container.NewMax(
		container.NewVBox(
			widget.NewLabelWithStyle("warehouses", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true}),
			container.NewGridWrap(fyne.NewSize(600, 400), table),
		),
	)

	createButton := widget.NewButton("Create", func() {
		m.ShowCreateWarehouseDialog(window)
	})

	updateButton := widget.NewButton("Update", func() {
		m.ShowUpdateWarehouseDialog(window)
	})

	deleteButton := widget.NewButton("Delete", func() {
		m.ShowDeleteWarehouseDialog(window)
	})

	exitButton := widget.NewButton("Back", func() {
		m.ShowMainScreen(window, m.UserLabel.Text)
	})

	topButtons := container.NewHBox(
		widget.NewSeparator(),
		container.NewGridWithColumns(1, exitButton),
	)

	buttons := container.NewHBox(
		widget.NewSeparator(),
		container.NewGridWithColumns(3, createButton, updateButton, deleteButton),
	)

	content := container.NewBorder(
		topButtons,
		buttons,
		nil,
		nil,
		tableContainer,
	)

	window.SetContent(content)
}

// ShowCreateWarehouseDialog shows user's form for warehouse's records creation
func (m *AppManager) ShowCreateWarehouseDialog(window fyne.Window) {
	nameEntry := widget.NewEntry()
	quantityEntry := widget.NewEntry()
	amountEntry := widget.NewEntry()

	dialog.ShowForm("Create Warehouse's record", "Create", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("name", nameEntry),
			widget.NewFormItem("quantity", quantityEntry),
			widget.NewFormItem("amount", amountEntry),
		}, func(confirmed bool) {
			if confirmed {
				quantity, err := strconv.Atoi(quantityEntry.Text)
				if err == nil {
					fmt.Printf("cannot convert text quantity to integer")
				}
				amount, err := strconv.Atoi(amountEntry.Text)
				if err == nil {
					fmt.Printf("cannot convert text amount to integer")
				}

				err = m.ShopService.CreateWarehousesItem(context.Background(), &dto.WarehousesData{
					Name:     nameEntry.Text,
					Quantity: quantity,
					Amount:   amount,
				})
				if err != nil {
					dialog.ShowError(err, window)
				} else {
					m.ShowWarehousesTable(window)
				}
			}
		}, window)
}

// ShowUpdateWarehouseDialog shows user's form for warehouse's records update
func (m *AppManager) ShowUpdateWarehouseDialog(window fyne.Window) {
	idEntry := widget.NewEntry()
	nameEntry := widget.NewEntry()
	quantityEntry := widget.NewEntry()
	amountEntry := widget.NewEntry()

	dialog.ShowForm("Please, enter Id", "Send", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("id", idEntry),
		}, func(confirmed bool) {
			if confirmed {
				dialog.ShowForm("Update Warehouse's record", "Update", "Cancel",
					[]*widget.FormItem{
						widget.NewFormItem("name", nameEntry),
						widget.NewFormItem("quantity", quantityEntry),
						widget.NewFormItem("amount", amountEntry),
					}, func(confirmed bool) {
						if confirmed {
							id, err := strconv.Atoi(idEntry.Text)
							if err == nil {
								fmt.Printf("cannot convert text id to integer")
							}
							quantity, err := strconv.Atoi(quantityEntry.Text)
							if err == nil {
								fmt.Printf("cannot convert text quantity to integer")
							}
							amount, err := strconv.Atoi(amountEntry.Text)
							if err == nil {
								fmt.Printf("cannot convert text amount to integer")
							}

							err = m.ShopService.UpdateWarehousesItem(context.Background(), &dto.WarehousesData{
								Id:       id,
								Name:     nameEntry.Text,
								Quantity: quantity,
								Amount:   amount,
							})
							if err != nil {
								dialog.ShowError(err, window)
							} else {
								m.ShowWarehousesTable(window)
							}
						}
					}, window)
			}
		}, window)
}

// ShowDeleteWarehouseDialog shows user's form for warehouse's records deleting
func (m *AppManager) ShowDeleteWarehouseDialog(window fyne.Window) {
	idEntry := widget.NewEntry()

	dialog.ShowForm("Delete Warehouse's record", "Delete", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Please, enter id", idEntry),
		}, func(confirmed bool) {
			if confirmed {
				id, err := strconv.Atoi(idEntry.Text)
				if err == nil {
					fmt.Printf("cannot convert text id to integer")
				}

				err = m.ShopService.DeleteWarehousesItem(context.Background(), id)
				if err != nil {
					dialog.ShowError(err, window)
				} else {
					m.ShowWarehousesTable(window)
				}
			}
		}, window)
}

// ShowExpenseItemsTable outputs data from expense items table
func (m *AppManager) ShowExpenseItemsTable(window fyne.Window) {
	data, err := m.ShopService.ShowExpenseItemsTable(context.Background())
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	headers := []string{"id", "name"}

	table := widget.NewTable(
		func() (int, int) { return len(data) + 1, len(headers) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			if id.Row == 0 {
				label.SetText(headers[id.Col])
				label.TextStyle = fyne.TextStyle{Monospace: true}
				return
			}

			row := id.Row - 1
			switch id.Col {
			case 0:
				label.SetText(strconv.Itoa(data[row].Id))
			case 1:
				label.SetText(data[row].Name)
			}
			label.TextStyle = fyne.TextStyle{Monospace: true}
		},
	)

	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 200)

	tableContainer := container.NewMax(
		container.NewVBox(
			widget.NewLabelWithStyle("expense_items", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true, Bold: true}),
			container.NewGridWrap(fyne.NewSize(600, 400), table),
		),
	)

	createButton := widget.NewButton("Create", func() {
		m.ShowCreateExpenseItemsDialog(window)
	})

	updateButton := widget.NewButton("Update", func() {
		m.ShowUpdateExpenseItemsDialog(window)
	})

	deleteButton := widget.NewButton("Delete", func() {
		m.ShowDeleteExpenseItemsDialog(window)
	})

	exitButton := widget.NewButton("Back", func() {
		m.ShowMainScreen(window, m.UserLabel.Text)
	})

	topButtons := container.NewHBox(
		widget.NewSeparator(),
		container.NewGridWithColumns(1, exitButton),
	)

	buttons := container.NewHBox(
		widget.NewSeparator(),
		container.NewGridWithColumns(3, createButton, updateButton, deleteButton),
	)

	content := container.NewBorder(
		topButtons,
		buttons,
		nil,
		nil,
		tableContainer,
	)

	window.SetContent(content)
}

// ShowCreateExpenseItemsDialog shows user's form for expense item's records creation
func (m *AppManager) ShowCreateExpenseItemsDialog(window fyne.Window) {
	nameEntry := widget.NewEntry()

	dialog.ShowForm("Create Expense Item's record", "Create", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("name", nameEntry),
		}, func(confirmed bool) {
			if confirmed {
				err := m.ShopService.CreateExpenseItem(context.Background(), nameEntry.Text)
				if err != nil {
					dialog.ShowError(err, window)
				} else {
					m.ShowExpenseItemsTable(window)
				}
			}
		}, window)
}

// ShowUpdateExpenseItemsDialog shows user's form for expense item's records update
func (m *AppManager) ShowUpdateExpenseItemsDialog(window fyne.Window) {
	idEntry := widget.NewEntry()
	nameEntry := widget.NewEntry()
	dialog.ShowForm("Please, enter Id", "Send", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("id", idEntry),
		}, func(confirmed bool) {
			if confirmed {
				dialog.ShowForm("Update Expense Item's record", "Update", "Cancel",
					[]*widget.FormItem{
						widget.NewFormItem("name", nameEntry),
					}, func(confirmed bool) {
						if confirmed {
							id, err := strconv.Atoi(idEntry.Text)
							if err == nil {
								fmt.Printf("cannot convert text id to integer")
							}

							err = m.ShopService.UpdateExpenseItem(context.Background(), &dto.ExpenseItemsData{
								Id:   id,
								Name: nameEntry.Text,
							})
							if err != nil {
								dialog.ShowError(err, window)
							} else {
								m.ShowExpenseItemsTable(window)
							}
						}
					}, window)
			}
		}, window)
}

// ShowDeleteExpenseItemsDialog shows user's form for expense item's records deleting
func (m *AppManager) ShowDeleteExpenseItemsDialog(window fyne.Window) {
	idEntry := widget.NewEntry()

	dialog.ShowForm("Delete Expense Item's record", "Delete", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Please, enter id", idEntry),
		}, func(confirmed bool) {
			if confirmed {
				id, err := strconv.Atoi(idEntry.Text)
				if err == nil {
					fmt.Printf("cannot convert text id to integer")
				}

				err = m.ShopService.DeleteExpenseItem(context.Background(), id)
				if err != nil {
					dialog.ShowError(err, window)
				} else {
					m.ShowExpenseItemsTable(window)
				}
			}
		}, window)
}

// ShowJournalsScreen shows screen with journals' features to user
func (m *AppManager) ShowJournalsScreen(window fyne.Window) fyne.CanvasObject {
	chargesButton := widget.NewButton("Charges", func() {
		m.ShowChargesTable(window)
	})

	salesButton := widget.NewButton("Sales", func() {
		m.ShowSalesTable(window)
	})

	return container.NewVBox(
		widget.NewLabelWithStyle("Please, choose Journal:", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}),
		chargesButton,
		salesButton,
	)
}

// ShowChargesTable outputs data from charges table
func (m *AppManager) ShowChargesTable(window fyne.Window) {
	data, err := m.ShopService.ShowChargesTable(context.Background())
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	headers := []string{"id", "amount", "charge_date", "expense_item_id"}

	table := widget.NewTable(
		func() (int, int) { return len(data) + 1, len(headers) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			if id.Row == 0 {
				label.SetText(headers[id.Col])
				label.TextStyle = fyne.TextStyle{Monospace: true}
				return
			}

			row := id.Row - 1
			switch id.Col {
			case 0:
				label.SetText(strconv.Itoa(data[row].Id))
			case 1:
				label.SetText(strconv.Itoa(data[row].Amount))
			case 2:
				label.SetText(data[row].ChargeDate)
			case 3:
				label.SetText(strconv.Itoa(data[row].ExpenseItemId))
			}
			label.TextStyle = fyne.TextStyle{Monospace: true}
		},
	)

	table.SetColumnWidth(0, 50)  // ID
	table.SetColumnWidth(1, 200) // Amount
	table.SetColumnWidth(2, 200) // Charge date
	table.SetColumnWidth(3, 50)  // Expense item id

	tableContainer := container.NewMax(
		container.NewVBox(
			widget.NewLabelWithStyle("charges", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true}),
			container.NewGridWrap(fyne.NewSize(600, 400), table),
		),
	)

	createButton := widget.NewButton("Create", func() {
		m.ShowCreateChargesDialog(window)
	})

	updateButton := widget.NewButton("Update", func() {
		m.ShowUpdateChargesDialog(window)
	})

	deleteButton := widget.NewButton("Delete", func() {
		m.ShowDeleteChargesDialog(window)
	})

	exitButton := widget.NewButton("Back", func() {
		m.ShowMainScreen(window, m.UserLabel.Text)
	})

	topButtons := container.NewHBox(
		widget.NewSeparator(),
		container.NewGridWithColumns(1, exitButton),
	)

	buttons := container.NewHBox(
		widget.NewSeparator(),
		container.NewGridWithColumns(3, createButton, updateButton, deleteButton),
	)

	content := container.NewBorder(
		topButtons,
		buttons,
		nil,
		nil,
		tableContainer,
	)

	window.SetContent(content)
}

// ShowCreateChargesDialog shows user's form for charges' records creation
func (m *AppManager) ShowCreateChargesDialog(window fyne.Window) {
	amountEntry := widget.NewEntry()
	chargeDateEntry := widget.NewEntry()
	expenseItemIdEntry := widget.NewEntry()

	dialog.ShowForm("Create Charges' record", "Create", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("amount", amountEntry),
			widget.NewFormItem("charge date", chargeDateEntry),
			widget.NewFormItem("expense item id", expenseItemIdEntry),
		}, func(confirmed bool) {
			if confirmed {
				exItemId, err := strconv.Atoi(expenseItemIdEntry.Text)
				if err == nil {
					fmt.Printf("cannot convert text expense item id to integer")
				}
				amount, err := strconv.Atoi(amountEntry.Text)
				if err == nil {
					fmt.Printf("cannot convert text amount to integer")
				}

				err = m.ShopService.CreateChargesItem(context.Background(), &dto.ChargesData{
					Amount:        amount,
					ChargeDate:    chargeDateEntry.Text,
					ExpenseItemId: exItemId,
				})
				if err != nil {
					dialog.ShowError(err, window)
				} else {
					m.ShowChargesTable(window)
				}
			}
		}, window)
}

// ShowUpdateChargesDialog shows user's form for charges' records update
func (m *AppManager) ShowUpdateChargesDialog(window fyne.Window) {
	idEntry := widget.NewEntry()
	amountEntry := widget.NewEntry()
	chargeDateEntry := widget.NewEntry()
	expenseItemIdEntry := widget.NewEntry()

	dialog.ShowForm("Please, enter Id", "Send", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("id", idEntry),
		}, func(confirmed bool) {
			if confirmed {
				dialog.ShowForm("Update Charges' record", "Update", "Cancel",
					[]*widget.FormItem{
						widget.NewFormItem("amount", amountEntry),
						widget.NewFormItem("charge date", chargeDateEntry),
						widget.NewFormItem("expense item id", expenseItemIdEntry),
					}, func(confirmed bool) {
						if confirmed {
							id, err := strconv.Atoi(idEntry.Text)
							if err == nil {
								fmt.Printf("cannot convert text id to integer")
							}
							exItemId, err := strconv.Atoi(expenseItemIdEntry.Text)
							if err == nil {
								fmt.Printf("cannot convert text expense item id to integer")
							}
							amount, err := strconv.Atoi(amountEntry.Text)
							if err == nil {
								fmt.Printf("cannot convert text amount to integer")
							}

							err = m.ShopService.UpdateChargesItem(context.Background(), &dto.ChargesData{
								Id:            id,
								Amount:        amount,
								ChargeDate:    chargeDateEntry.Text,
								ExpenseItemId: exItemId,
							})
							if err != nil {
								dialog.ShowError(err, window)
							} else {
								m.ShowChargesTable(window)
							}
						}
					}, window)
			}
		}, window)
}

// ShowDeleteChargesDialog shows user's form for charges' records deleting
func (m *AppManager) ShowDeleteChargesDialog(window fyne.Window) {
	idEntry := widget.NewEntry()

	dialog.ShowForm("Delete Charges' record", "Delete", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Please, enter id", idEntry),
		}, func(confirmed bool) {
			if confirmed {
				id, err := strconv.Atoi(idEntry.Text)
				if err == nil {
					fmt.Printf("cannot convert text id to integer")
				}

				err = m.ShopService.DeleteChargesItem(context.Background(), id)
				if err != nil {
					dialog.ShowError(err, window)
				} else {
					m.ShowChargesTable(window)
				}
			}
		}, window)
}

// ShowSalesTable outputs data from sales table
func (m *AppManager) ShowSalesTable(window fyne.Window) {
	data, err := m.ShopService.ShowSalesTable(context.Background())
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	headers := []string{"id", "amount", "quantity", "sale_date", "warehouses_id"}

	table := widget.NewTable(
		func() (int, int) { return len(data) + 1, len(headers) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			if id.Row == 0 {
				label.SetText(headers[id.Col])
				label.TextStyle = fyne.TextStyle{Monospace: true}
				return
			}

			row := id.Row - 1
			switch id.Col {
			case 0:
				label.SetText(strconv.Itoa(data[row].Id))
			case 1:
				label.SetText(strconv.Itoa(data[row].Amount))
			case 2:
				label.SetText(strconv.Itoa(data[row].Quantity))
			case 3:
				label.SetText(data[row].SaleDate)
			case 4:
				label.SetText(strconv.Itoa(data[row].WarehousesId))
			}
			label.TextStyle = fyne.TextStyle{Monospace: true}
		},
	)

	table.SetColumnWidth(0, 50)  // ID
	table.SetColumnWidth(1, 100) // Amount
	table.SetColumnWidth(2, 100) // Quantity
	table.SetColumnWidth(3, 200) // Sale date
	table.SetColumnWidth(4, 50)  // Warehouses id

	tableContainer := container.NewMax(
		container.NewVBox(
			widget.NewLabelWithStyle("sales", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true}),
			container.NewGridWrap(fyne.NewSize(600, 400), table),
		),
	)

	createButton := widget.NewButton("Create", func() {
		m.ShowCreateSalesDialog(window)
	})

	updateButton := widget.NewButton("Update", func() {
		m.ShowUpdateSalesDialog(window)
	})

	deleteButton := widget.NewButton("Delete", func() {
		m.ShowDeleteSalesDialog(window)
	})

	exitButton := widget.NewButton("Back", func() {
		m.ShowMainScreen(window, m.UserLabel.Text)
	})

	topButtons := container.NewHBox(
		widget.NewSeparator(),
		container.NewGridWithColumns(1, exitButton),
	)

	buttons := container.NewHBox(
		widget.NewSeparator(),
		container.NewGridWithColumns(3, createButton, updateButton, deleteButton),
	)

	content := container.NewBorder(
		topButtons,
		buttons,
		nil,
		nil,
		tableContainer,
	)

	window.SetContent(content)
}

// ShowCreateSalesDialog shows user's form for sales' records creation
func (m *AppManager) ShowCreateSalesDialog(window fyne.Window) {
	amountEntry := widget.NewEntry()
	quantityEntry := widget.NewEntry()
	saleDateEntry := widget.NewEntry()
	warehousesIdEntry := widget.NewEntry()

	dialog.ShowForm("Create Sales' record", "Create", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("amount", amountEntry),
			widget.NewFormItem("quantity", quantityEntry),
			widget.NewFormItem("sale date", saleDateEntry),
			widget.NewFormItem("warehouses id", warehousesIdEntry),
		}, func(confirmed bool) {
			if confirmed {
				warehousesId, err := strconv.Atoi(warehousesIdEntry.Text)
				if err == nil {
					fmt.Printf("cannot convert text warehouses id to integer")
				}
				amount, err := strconv.Atoi(amountEntry.Text)
				if err == nil {
					fmt.Printf("cannot convert text amount to integer")
				}
				quantity, err := strconv.Atoi(quantityEntry.Text)
				if err == nil {
					fmt.Printf("cannot convert text quantity to integer")
				}

				err = m.ShopService.CreateSalesItem(context.Background(), &dto.SalesData{
					Amount:       amount,
					Quantity:     quantity,
					SaleDate:     saleDateEntry.Text,
					WarehousesId: warehousesId,
				})
				if err != nil {
					dialog.ShowError(err, window)
				} else {
					m.ShowSalesTable(window)
				}
			}
		}, window)
}

// ShowUpdateSalesDialog shows user's form for sales' records update
func (m *AppManager) ShowUpdateSalesDialog(window fyne.Window) {
	amountEntry := widget.NewEntry()
	quantityEntry := widget.NewEntry()
	saleDateEntry := widget.NewEntry()
	warehousesIdEntry := widget.NewEntry()
	idEntry := widget.NewEntry()

	dialog.ShowForm("Please, enter Id", "Send", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("id", idEntry),
		}, func(confirmed bool) {
			if confirmed {
				dialog.ShowForm("Update Sales' record", "Update", "Cancel",
					[]*widget.FormItem{
						widget.NewFormItem("amount", amountEntry),
						widget.NewFormItem("quantity", quantityEntry),
						widget.NewFormItem("sale date", saleDateEntry),
						widget.NewFormItem("warehouses id", warehousesIdEntry),
					}, func(confirmed bool) {
						if confirmed {
							id, err := strconv.Atoi(idEntry.Text)
							if err == nil {
								fmt.Printf("cannot convert text id to integer")
							}
							warehousesId, err := strconv.Atoi(warehousesIdEntry.Text)
							if err == nil {
								fmt.Printf("cannot convert text warehouses id to integer")
							}
							amount, err := strconv.Atoi(amountEntry.Text)
							if err == nil {
								fmt.Printf("cannot convert text amount to integer")
							}
							quantity, err := strconv.Atoi(quantityEntry.Text)
							if err == nil {
								fmt.Printf("cannot convert text quantity to integer")
							}

							err = m.ShopService.UpdateSalesItem(context.Background(), &dto.SalesData{
								Id:           id,
								Amount:       amount,
								Quantity:     quantity,
								SaleDate:     saleDateEntry.Text,
								WarehousesId: warehousesId,
							})
							if err != nil {
								dialog.ShowError(err, window)
							} else {
								m.ShowSalesTable(window)
							}
						}
					}, window)
			}
		}, window)
}

// ShowDeleteSalesDialog shows user's form for sales' records update
func (m *AppManager) ShowDeleteSalesDialog(window fyne.Window) {
	idEntry := widget.NewEntry()

	dialog.ShowForm("Delete Sales' record", "Delete", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Please, enter id", idEntry),
		}, func(confirmed bool) {
			if confirmed {
				id, err := strconv.Atoi(idEntry.Text)
				if err == nil {
					fmt.Printf("cannot convert text id to integer")
				}

				err = m.ShopService.DeleteSalesItem(context.Background(), id)
				if err != nil {
					dialog.ShowError(err, window)
				} else {
					m.ShowSalesTable(window)
				}
			}
		}, window)
}

// ShowReportsScreen shows screen with reports' features to user
func (m *AppManager) ShowReportsScreen(window fyne.Window) fyne.CanvasObject {
	profitButton := widget.NewButton("Count month profit", func() {
		m.ShowMonthProfit(window)
	})

	itemsButton := widget.NewButton("Show 5 best items", func() {
		m.ShowBestItems(window)
	})

	return container.NewVBox(
		widget.NewLabelWithStyle("Please, choose Report:", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}),
		profitButton,
		itemsButton,
	)
}

// ShowMonthProfit counts summary month profit of shop and outputs it
func (m *AppManager) ShowMonthProfit(window fyne.Window) {
	profit, err := m.ShopService.CountMonthProfit(context.Background())
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to count monthly profit: %v", err), window)
		return
	}

	message := fmt.Sprintf("Summary month profit: %d", profit)

	dialog.ShowInformation("Profit", message, window)
}

// ShowBestItems outputs additional table which contains data of five most profitable items in shop
func (m *AppManager) ShowBestItems(window fyne.Window) {
	fromEntry := widget.NewEntry()
	toEntry := widget.NewEntry()

	dialog.ShowForm("Please, enter date range", "Approve", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("from", fromEntry),
			widget.NewFormItem("to", toEntry),
		}, func(confirmed bool) {
			if confirmed {
				data, err := m.ShopService.GetFiveBestItems(context.Background(), fromEntry.Text, toEntry.Text)
				if err != nil {
					dialog.ShowError(err, window)
					return
				}

				headers := []string{"name", "total_revenue"}

				table := widget.NewTable(
					func() (int, int) { return len(data) + 1, len(headers) },
					func() fyne.CanvasObject { return widget.NewLabel("") },
					func(id widget.TableCellID, cell fyne.CanvasObject) {
						label := cell.(*widget.Label)

						if id.Row == 0 {
							label.SetText(headers[id.Col])
							label.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
							return
						}

						row := id.Row - 1
						switch id.Col {
						case 0:
							label.SetText(data[row].Name)
						case 1:
							label.SetText(strconv.Itoa(data[row].TotalRevenue))
						}
						label.TextStyle = fyne.TextStyle{Monospace: true}
					},
				)

				table.SetColumnWidth(0, 300)
				table.SetColumnWidth(1, 100)

				tableContainer := container.NewMax(
					container.NewVBox(
						widget.NewLabelWithStyle("five_best_items", fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true}),
						container.NewGridWrap(fyne.NewSize(600, 400), table),
					),
				)

				downloadButton := widget.NewButton("Download PDF", func() {
					m.generateBestItemsPDF(data, fromEntry.Text, toEntry.Text, window)
				})

				exitButton := widget.NewButton("Back", func() {
					m.ShowMainScreen(window, m.UserLabel.Text)
				})

				buttons := container.NewHBox(downloadButton, exitButton)

				content := container.NewBorder(nil, buttons, nil, nil, tableContainer)
				window.SetContent(content)
			}
		}, window)
}

// generateBestItemsPDF creates .pdf file with FiveBestItems report in root dir
func (m *AppManager) generateBestItemsPDF(data []*dto.BestItemsData, fromDate, toDate string, window fyne.Window) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Report: Five Best Items")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, "Date Range: "+fromDate+" to "+toDate)
	pdf.Ln(12)

	pdf.SetFont("Monospace", "B", 12)
	pdf.Cell(120, 10, "name")
	pdf.Cell(0, 10, "total_Revenue")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	for _, item := range data {
		pdf.Cell(120, 10, item.Name)
		pdf.Cell(0, 10, strconv.Itoa(item.TotalRevenue))
		pdf.Ln(8)
	}

	reportDir := "reports"
	err := os.MkdirAll(reportDir, os.ModePerm)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to create directory: %w", err), window)
		return
	}

	outputPath := filepath.Join(reportDir, "BestItemsReport.pdf")
	err = pdf.OutputFileAndClose(outputPath)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	dialog.ShowInformation("Download Complete", "Report saved to: "+outputPath, window)
}
