package graphics

import (
	"automatedShop/internal/services"
	"automatedShop/internal/services/dto"
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jung-kurt/gofpdf"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
)

type AppManager struct {
	AuthService services.IAuthService
	ShopService services.IShopService
}

func NewAppManager(s *services.Service) *AppManager {
	return &AppManager{
		AuthService: s.AuthService,
		ShopService: s.ShopService,
	}
}

func (m *AppManager) Run() {
	application := app.New()
	application.Settings().SetTheme(&CustomTheme{})
	mainWindow := application.NewWindow("Shop Management System v.0.0.0")
	mainWindow.Resize(fyne.NewSize(300, 300))

	m.ShowLoginScreen(mainWindow)

	mainWindow.ShowAndRun()
}

// ShowLoginScreen shows login screen window to user
func (m *AppManager) ShowLoginScreen(window fyne.Window) {
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Username")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	errorLabel := widget.NewLabel("")

	background := canvas.NewRectangle(color.RGBA{R: 220, G: 220, B: 255, A: 255})

	loginButton := NewSquareButton("Login", func() {
		username := usernameEntry.Text
		password := passwordEntry.Text

		if m.AuthService.AuthoriseUser(context.Background(), username, password) {
			m.showMainScreen(window)
		} else {
			errorLabel.SetText("Invalid username or password")
		}
	})

	registerButton := widget.NewButton("Register", func() {
		m.showRegisterScreen(window)
	})

	window.SetContent(container.NewStack(
		background,
		container.NewVBox(
			widget.NewLabel("Authorization"),
			usernameEntry,
			passwordEntry,
			loginButton,
			registerButton,
			errorLabel,
		)))
}

// Register Screen
func (m *AppManager) showRegisterScreen(window fyne.Window) {
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

	form := container.NewVBox(
		widget.NewLabel("Register"),
		usernameEntry,
		passwordEntry,
		registerButton,
		errorLabel,
	)

	window.SetContent(form)
}

// Main Screen
func (m *AppManager) showMainScreen(window fyne.Window) {
	tabs := container.NewAppTabs(
		container.NewTabItem("Handbooks", m.showHandbooksScreen(window)),
		container.NewTabItem("Journals", m.showJournalsScreen(window)),
		container.NewTabItem("Reports", m.showReportsScreen(window)),
	)

	exitButton := widget.NewButton("Logout", func() {
		m.ShowLoginScreen(window)
	})

	mainLayout := container.NewBorder(nil, exitButton, nil, nil, tabs)
	window.SetContent(mainLayout)
}

// Handbooks Screen
func (m *AppManager) showHandbooksScreen(window fyne.Window) fyne.CanvasObject {
	warehousesButton := widget.NewButton("Warehouses", func() {
		m.showWarehousesTable(window)
	})

	expenseItemsButton := widget.NewButton("Expense Items", func() {
		m.showExpenseItemsTable(window)
	})

	return container.NewVBox(
		widget.NewLabel("Choose Handbook:"),
		warehousesButton,
		expenseItemsButton,
	)
}

func (m *AppManager) showWarehousesTable(window fyne.Window) {
	data, err := m.ShopService.ShowWarehousesTable(context.Background())
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	// Заголовки столбцов
	headers := []string{"ID", "Name", "Quantity", "Amount"}

	// Таблица с данными
	table := widget.NewTable(
		func() (int, int) { return len(data) + 1, len(headers) }, // +1 для строки заголовков
		func() fyne.CanvasObject { return widget.NewLabel("") },  // Создаем ячейку
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			// Первая строка - заголовки
			if id.Row == 0 {
				label.SetText(headers[id.Col])
				label.TextStyle = fyne.TextStyle{Bold: true} // Жирный текст для заголовков
				return
			}

			// Остальные строки - данные
			row := id.Row - 1 // Убираем строку заголовка
			switch id.Col {
			case 0:
				label.SetText(strconv.Itoa(data[row].Id)) // ID
			case 1:
				label.SetText(data[row].Name) // Name
			case 2:
				label.SetText(strconv.Itoa(data[row].Quantity)) // Quantity
			case 3:
				label.SetText(strconv.Itoa(data[row].Amount)) // Amount
			}
		},
	)

	// Устанавливаем ширину колонок для улучшения читаемости
	table.SetColumnWidth(0, 50)  // ID
	table.SetColumnWidth(1, 200) // Name
	table.SetColumnWidth(2, 100) // Quantity
	table.SetColumnWidth(3, 100) // Amount

	// Задаем размеры таблицы, чтобы она занимала 60% окна
	tableContainer := container.NewMax(
		container.NewVBox(
			widget.NewLabelWithStyle("Warehouses", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			container.NewGridWrap(fyne.NewSize(600, 400), table), // 60% от окна (примерный размер)
		),
	)

	// Кнопки управления
	createButton := widget.NewButton("Create", func() {
		m.showCreateWarehouseDialog(window)
	})

	updateButton := widget.NewButton("Update", func() {
		m.showUpdateWarehouseDialog(window)
	})

	deleteButton := widget.NewButton("Delete", func() {
		m.showDeleteWarehouseDialog(window)
	})

	exitButton := widget.NewButton("Back", func() {
		m.showMainScreen(window)
	})

	// Расположение кнопок внизу
	topButtons := container.NewHBox(
		widget.NewSeparator(),
		container.NewGridWithColumns(1, exitButton),
	)

	buttons := container.NewHBox(
		widget.NewSeparator(),
		container.NewGridWithColumns(3, createButton, updateButton, deleteButton),
	)

	// Основной контент с таблицей (в центре) и кнопками (внизу)
	content := container.NewBorder(
		topButtons,     // Верхняя граница
		buttons,        // Нижняя граница
		nil,            // Левая граница
		nil,            // Правая граница
		tableContainer, // Центральная часть
	)

	// Устанавливаем контент окна
	window.SetContent(content)
}

// Create Warehouse Dialog
func (m *AppManager) showCreateWarehouseDialog(window fyne.Window) {
	nameEntry := widget.NewEntry()
	quantityEntry := widget.NewEntry()
	amountEntry := widget.NewEntry()

	dialog.ShowForm("Create Warehouse", "Create", "Cancel",
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
					m.showWarehousesTable(window)
				}
			}
		}, window)
}

// Update Warehouse Dialog
func (m *AppManager) showUpdateWarehouseDialog(window fyne.Window) {
	idEntry := widget.NewEntry()
	nameEntry := widget.NewEntry()
	quantityEntry := widget.NewEntry()
	amountEntry := widget.NewEntry()
	dialog.ShowForm("Please, enter Id", "Send", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("id", idEntry),
		}, func(confirmed bool) {
			if confirmed {
				dialog.ShowForm("Update Warehouse", "Update", "Cancel",
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
								m.showWarehousesTable(window)
							}
						}
					}, window)
			}
		}, window)
}

// Delete Warehouse Dialog
func (m *AppManager) showDeleteWarehouseDialog(window fyne.Window) {
	idEntry := widget.NewEntry()

	dialog.ShowForm("Delete Warehouse", "Delete", "Cancel",
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
					m.showWarehousesTable(window)
				}
			}
		}, window)
}

// Expense Items Table
func (m *AppManager) showExpenseItemsTable(window fyne.Window) {
	data, err := m.ShopService.ShowExpenseItemsTable(context.Background())
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	// Заголовки столбцов
	headers := []string{"ID", "Name"}

	// Таблица с данными
	table := widget.NewTable(
		func() (int, int) { return len(data) + 1, len(headers) }, // +1 для строки заголовков
		func() fyne.CanvasObject { return widget.NewLabel("") },  // Создаем ячейку
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			// Первая строка - заголовки
			if id.Row == 0 {
				label.SetText(headers[id.Col])
				label.TextStyle = fyne.TextStyle{Bold: true} // Жирный текст для заголовков
				return
			}

			// Остальные строки - данные
			row := id.Row - 1 // Убираем строку заголовка
			switch id.Col {
			case 0:
				label.SetText(strconv.Itoa(data[row].Id)) // ID
			case 1:
				label.SetText(data[row].Name) // Name
			}
		},
	)

	// Устанавливаем ширину колонок для улучшения читаемости
	table.SetColumnWidth(0, 50)  // ID
	table.SetColumnWidth(1, 200) // Name

	// Задаем размеры таблицы, чтобы она занимала 60% окна
	tableContainer := container.NewMax(
		container.NewVBox(
			widget.NewLabelWithStyle("Expense Items", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			container.NewGridWrap(fyne.NewSize(600, 400), table), // 60% от окна (примерный размер)
		),
	)

	// Кнопки управления
	createButton := widget.NewButton("Create", func() {
		m.showCreateExpenseItemsDialog(window)
	})

	updateButton := widget.NewButton("Update", func() {
		m.showUpdateExpenseItemsDialog(window)
	})

	deleteButton := widget.NewButton("Delete", func() {
		m.showDeleteExpenseItemsDialog(window)
	})

	exitButton := widget.NewButton("Back", func() {
		m.showMainScreen(window)
	})

	// Расположение кнопок внизу
	topButtons := container.NewHBox(
		widget.NewSeparator(),
		container.NewGridWithColumns(1, exitButton),
	)

	buttons := container.NewHBox(
		widget.NewSeparator(),
		container.NewGridWithColumns(3, createButton, updateButton, deleteButton),
	)

	// Основной контент с таблицей (в центре) и кнопками (внизу)
	content := container.NewBorder(
		topButtons,     // Верхняя граница
		buttons,        // Нижняя граница
		nil,            // Левая граница
		nil,            // Правая граница
		tableContainer, // Центральная часть
	)

	// Устанавливаем контент окна
	window.SetContent(content)
}

// Create Expense Items Dialog
func (m *AppManager) showCreateExpenseItemsDialog(window fyne.Window) {
	nameEntry := widget.NewEntry()

	dialog.ShowForm("Create Expense Item", "Create", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("name", nameEntry),
		}, func(confirmed bool) {
			if confirmed {
				err := m.ShopService.CreateExpenseItem(context.Background(), nameEntry.Text)
				if err != nil {
					dialog.ShowError(err, window)
				} else {
					m.showExpenseItemsTable(window)
				}
			}
		}, window)
}

// Update Expense Item Dialog
func (m *AppManager) showUpdateExpenseItemsDialog(window fyne.Window) {
	idEntry := widget.NewEntry()
	nameEntry := widget.NewEntry()
	dialog.ShowForm("Please, enter Id", "Send", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("id", idEntry),
		}, func(confirmed bool) {
			if confirmed {
				dialog.ShowForm("Update Expense Item", "Update", "Cancel",
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
								m.showExpenseItemsTable(window)
							}
						}
					}, window)
			}
		}, window)
}

// Delete Expense Item Dialog
func (m *AppManager) showDeleteExpenseItemsDialog(window fyne.Window) {
	idEntry := widget.NewEntry()

	dialog.ShowForm("Delete Expense Item", "Delete", "Cancel",
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
					m.showExpenseItemsTable(window)
				}
			}
		}, window)
}

// Journals Screen
func (m *AppManager) showJournalsScreen(window fyne.Window) fyne.CanvasObject {
	chargesButton := widget.NewButton("Charges", func() {
		m.showChargesTable(window)
	})

	salesButton := widget.NewButton("Sales", func() {
		m.showSalesTable(window)
	})

	return container.NewVBox(
		widget.NewLabel("Choose Journals:"),
		chargesButton,
		salesButton,
	)
}

// Charges Table
func (m *AppManager) showChargesTable(window fyne.Window) {
	data, err := m.ShopService.ShowChargesTable(context.Background())
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	// Заголовки столбцов
	headers := []string{"ID", "Amount", "Charge date", "Expense item id"}

	// Таблица с данными
	table := widget.NewTable(
		func() (int, int) { return len(data) + 1, len(headers) }, // +1 для строки заголовков
		func() fyne.CanvasObject { return widget.NewLabel("") },  // Создаем ячейку
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			// Первая строка - заголовки
			if id.Row == 0 {
				label.SetText(headers[id.Col])
				label.TextStyle = fyne.TextStyle{Bold: true} // Жирный текст для заголовков
				return
			}

			// Остальные строки - данные
			row := id.Row - 1 // Убираем строку заголовка
			switch id.Col {
			case 0:
				label.SetText(strconv.Itoa(data[row].Id)) // ID
			case 1:
				label.SetText(strconv.Itoa(data[row].Amount)) // Amount
			case 2:
				label.SetText(data[row].ChargeDate) // Quantity
			case 3:
				label.SetText(strconv.Itoa(data[row].ExpenseItemId))
			}
		},
	)

	// Устанавливаем ширину колонок для улучшения читаемости
	table.SetColumnWidth(0, 50)  // ID
	table.SetColumnWidth(1, 200) // Amount
	table.SetColumnWidth(2, 200) // Charge date
	table.SetColumnWidth(3, 50)  // Expense item id

	// Задаем размеры таблицы, чтобы она занимала 60% окна
	tableContainer := container.NewMax(
		container.NewVBox(
			widget.NewLabelWithStyle("Charges", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			container.NewGridWrap(fyne.NewSize(600, 400), table), // 60% от окна (примерный размер)
		),
	)

	// Кнопки управления
	createButton := widget.NewButton("Create", func() {
		m.showCreateChargesDialog(window)
	})

	updateButton := widget.NewButton("Update", func() {
		m.showUpdateChargesDialog(window)
	})

	deleteButton := widget.NewButton("Delete", func() {
		m.showDeleteChargesDialog(window)
	})

	exitButton := widget.NewButton("Back", func() {
		m.showMainScreen(window)
	})

	// Расположение кнопок внизу
	topButtons := container.NewHBox(
		widget.NewSeparator(),
		container.NewGridWithColumns(1, exitButton),
	)

	buttons := container.NewHBox(
		widget.NewSeparator(),
		container.NewGridWithColumns(3, createButton, updateButton, deleteButton),
	)

	// Основной контент с таблицей (в центре) и кнопками (внизу)
	content := container.NewBorder(
		topButtons,     // Верхняя граница
		buttons,        // Нижняя граница
		nil,            // Левая граница
		nil,            // Правая граница
		tableContainer, // Центральная часть
	)

	// Устанавливаем контент окна
	window.SetContent(content)
}

// Create Charges Dialog
func (m *AppManager) showCreateChargesDialog(window fyne.Window) {
	amountEntry := widget.NewEntry()
	chargeDateEntry := widget.NewEntry()
	expenseItemIdEntry := widget.NewEntry()

	dialog.ShowForm("Create Charge", "Create", "Cancel",
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
					m.showChargesTable(window)
				}
			}
		}, window)
}

// Update Charges Dialog
func (m *AppManager) showUpdateChargesDialog(window fyne.Window) {
	idEntry := widget.NewEntry()
	amountEntry := widget.NewEntry()
	chargeDateEntry := widget.NewEntry()
	expenseItemIdEntry := widget.NewEntry()
	dialog.ShowForm("Please, enter Id", "Send", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("id", idEntry),
		}, func(confirmed bool) {
			if confirmed {
				dialog.ShowForm("UpdateCharges", "Update", "Cancel",
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
								m.showChargesTable(window)
							}
						}
					}, window)
			}
		}, window)
}

// Delete Charges Dialog
func (m *AppManager) showDeleteChargesDialog(window fyne.Window) {
	idEntry := widget.NewEntry()

	dialog.ShowForm("Delete Charges", "Delete", "Cancel",
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
					m.showChargesTable(window)
				}
			}
		}, window)
}

// Sales Table
func (m *AppManager) showSalesTable(window fyne.Window) {
	data, err := m.ShopService.ShowSalesTable(context.Background())
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	// Заголовки столбцов
	headers := []string{"ID", "Amount", "Quantity", "Sale date", "Warehouses id"}

	// Таблица с данными
	table := widget.NewTable(
		func() (int, int) { return len(data) + 1, len(headers) }, // +1 для строки заголовков
		func() fyne.CanvasObject { return widget.NewLabel("") },  // Создаем ячейку
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			// Первая строка - заголовки
			if id.Row == 0 {
				label.SetText(headers[id.Col])
				label.TextStyle = fyne.TextStyle{Bold: true} // Жирный текст для заголовков
				return
			}

			// Остальные строки - данные
			row := id.Row - 1 // Убираем строку заголовка
			switch id.Col {
			case 0:
				label.SetText(strconv.Itoa(data[row].Id)) // ID
			case 1:
				label.SetText(strconv.Itoa(data[row].Amount)) // Amount
			case 2:
				label.SetText(strconv.Itoa(data[row].Quantity)) // Quantity
			case 3:
				label.SetText(data[row].SaleDate)
			case 4:
				label.SetText(strconv.Itoa(data[row].WarehousesId))
			}
		},
	)

	// Устанавливаем ширину колонок для улучшения читаемости
	table.SetColumnWidth(0, 50)  // ID
	table.SetColumnWidth(1, 100) // Amount
	table.SetColumnWidth(2, 100) // Quantity
	table.SetColumnWidth(3, 200) // Sale date
	table.SetColumnWidth(4, 50)  // Warehouses id

	// Задаем размеры таблицы, чтобы она занимала 60% окна
	tableContainer := container.NewMax(
		container.NewVBox(
			widget.NewLabelWithStyle("Sales", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			container.NewGridWrap(fyne.NewSize(600, 400), table), // 60% от окна (примерный размер)
		),
	)

	// Кнопки управления
	createButton := widget.NewButton("Create", func() {
		m.showCreateSalesDialog(window)
	})

	updateButton := widget.NewButton("Update", func() {
		m.showUpdateSalesDialog(window)
	})

	deleteButton := widget.NewButton("Delete", func() {
		m.showDeleteSalesDialog(window)
	})

	exitButton := widget.NewButton("Back", func() {
		m.showMainScreen(window)
	})

	// Расположение кнопок внизу
	topButtons := container.NewHBox(
		widget.NewSeparator(),
		container.NewGridWithColumns(1, exitButton),
	)

	buttons := container.NewHBox(
		widget.NewSeparator(),
		container.NewGridWithColumns(3, createButton, updateButton, deleteButton),
	)

	// Основной контент с таблицей (в центре) и кнопками (внизу)
	content := container.NewBorder(
		topButtons,     // Верхняя граница
		buttons,        // Нижняя граница
		nil,            // Левая граница
		nil,            // Правая граница
		tableContainer, // Центральная часть
	)

	// Устанавливаем контент окна
	window.SetContent(content)
}

// Create Sales Dialog
func (m *AppManager) showCreateSalesDialog(window fyne.Window) {
	amountEntry := widget.NewEntry()
	quantityEntry := widget.NewEntry()
	saleDateEntry := widget.NewEntry()
	warehousesIdEntry := widget.NewEntry()

	dialog.ShowForm("Create Sale", "Create", "Cancel",
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
					m.showSalesTable(window)
				}
			}
		}, window)
}

// Update Sales Dialog
func (m *AppManager) showUpdateSalesDialog(window fyne.Window) {
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
				dialog.ShowForm("UpdateSales", "Update", "Cancel",
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
								m.showSalesTable(window)
							}
						}
					}, window)
			}
		}, window)
}

// Delete Sales Dialog
func (m *AppManager) showDeleteSalesDialog(window fyne.Window) {
	idEntry := widget.NewEntry()

	dialog.ShowForm("Delete Sales", "Delete", "Cancel",
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
					m.showSalesTable(window)
				}
			}
		}, window)
}

// Reports Screen
func (m *AppManager) showReportsScreen(window fyne.Window) fyne.CanvasObject {
	profitButton := widget.NewButton("Count month profit", func() {
		m.showMonthProfit(window)
	})

	itemsButton := widget.NewButton("Show 5 best items", func() {
		m.showBestItems(window)
	})

	return container.NewVBox(
		widget.NewLabel("Choose Report:"),
		profitButton,
		itemsButton,
	)
}

func (m *AppManager) showMonthProfit(window fyne.Window) {
	profit, err := m.ShopService.CountMonthProfit(context.Background())
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to count monthly profit: %v", err), window)
		return
	}

	// Формируем сообщение с прибылью
	message := fmt.Sprintf("Прибыль за текущий месяц: %d", profit)

	// Показываем информационное диалоговое окно с результатом
	dialog.ShowInformation("Monthly Profit", message, window)
}

func (m *AppManager) showBestItems(window fyne.Window) {
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

				headers := []string{"Name", "Total Revenue"}

				// Таблица с данными
				table := widget.NewTable(
					func() (int, int) { return len(data) + 1, len(headers) },
					func() fyne.CanvasObject { return widget.NewLabel("") },
					func(id widget.TableCellID, cell fyne.CanvasObject) {
						label := cell.(*widget.Label)

						if id.Row == 0 {
							label.SetText(headers[id.Col])
							label.TextStyle = fyne.TextStyle{Bold: true}
							return
						}

						row := id.Row - 1
						switch id.Col {
						case 0:
							label.SetText(data[row].Name)
						case 1:
							label.SetText(strconv.Itoa(data[row].TotalRevenue))
						}
					},
				)

				table.SetColumnWidth(0, 300)
				table.SetColumnWidth(1, 100)

				tableContainer := container.NewMax(
					container.NewVBox(
						widget.NewLabelWithStyle("Five Best Items", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
						container.NewGridWrap(fyne.NewSize(600, 400), table),
					),
				)

				downloadButton := widget.NewButton("Download PDF", func() {
					m.generateBestItemsPDF(data, fromEntry.Text, toEntry.Text, window)
				})

				exitButton := widget.NewButton("Back", func() {
					m.showMainScreen(window)
				})

				// Нижняя панель с кнопками
				buttons := container.NewHBox(downloadButton, exitButton)

				content := container.NewBorder(nil, buttons, nil, nil, tableContainer)
				window.SetContent(content)
			}
		}, window)
}

// Метод для генерации PDF-файла
func (m *AppManager) generateBestItemsPDF(data []*dto.BestItemsData, fromDate, toDate string, window fyne.Window) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Report: Five Best Items")
	pdf.Ln(12)

	// Печатаем дату
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, "Date Range: "+fromDate+" to "+toDate)
	pdf.Ln(12)

	// Заголовки таблицы
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(120, 10, "Name")
	pdf.Cell(0, 10, "Total Revenue")
	pdf.Ln(10)

	// Данные таблицы
	pdf.SetFont("Arial", "", 12)
	for _, item := range data {
		pdf.Cell(120, 10, item.Name)
		pdf.Cell(0, 10, strconv.Itoa(item.TotalRevenue))
		pdf.Ln(8)
	}

	// Создаём папку reports
	reportDir := "reports"
	err := os.MkdirAll(reportDir, os.ModePerm) // Создаёт папку, если её нет
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to create directory: %w", err), window)
		return
	}
	// Сохраняем PDF в файл
	outputPath := filepath.Join(reportDir, "BestItemsReport.pdf")
	err = pdf.OutputFileAndClose(outputPath)
	if err != nil {
		dialog.ShowError(err, window)
		return
	}

	// Показываем сообщение о сохранении
	dialog.ShowInformation("Download Complete", "Report saved to: "+outputPath, window)
}
