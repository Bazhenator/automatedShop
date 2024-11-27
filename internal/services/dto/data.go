package dto

type WarehousesData struct {
	Id       int
	Name     string
	Quantity int
	Amount   int
}

type ExpenseItemsData struct {
	Id   int
	Name string
}

type BestItemsData struct {
	Name         string
	TotalRevenue int
}

type SalesData struct {
	Id           int
	Amount       int
	Quantity     int
	SaleDate     string
	WarehousesId int
}

type ChargesData struct {
	Id            int
	Amount        int
	ChargeDate    string
	ExpenseItemId int
}
