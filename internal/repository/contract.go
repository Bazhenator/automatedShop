package repository

import (
	"automatedShop/internal/repository/dto"
	logicDto "automatedShop/internal/services/dto"
	"context"
)

type IShopRepository interface {
	// Handbook's methods
	ShowWarehousesTable(context.Context) ([]*logicDto.WarehousesData, error)
	CreateWarehousesItem(context.Context, string, int, int) error
	UpdateWarehousesItem(context.Context, string, int, int, int) error
	DeleteWarehousesItem(context.Context, int) error
	ShowExpenseItemsTable(context.Context) ([]*logicDto.ExpenseItemsData, error)
	CreateExpenseItem(context.Context, string) error
	UpdateExpenseItem(context.Context, string, int) error
	DeleteExpenseItem(context.Context, int) error

	// Journal's methods
	ShowChargesTable(context.Context) ([]*logicDto.ChargesData, error)
	CreateChargesItem(context.Context, *logicDto.ChargesData) error
	UpdateChargesItem(context.Context, *logicDto.ChargesData) error
	DeleteChargesItem(context.Context, int) error
	ShowSalesTable(context.Context) ([]*logicDto.SalesData, error)
	CreateSalesItem(context.Context, *logicDto.SalesData) error
	UpdateSalesItem(context.Context, *logicDto.SalesData) error
	DeleteSalesItem(context.Context, int) error

	// Report's methods
	CountMonthProfit(context.Context) (int64, error)
	GetFiveBestItems(context.Context, string, string) ([]*logicDto.BestItemsData, error)
}

type IAuthRepository interface {
	SaveUser(context.Context, string, []byte) error
	FindUser(context.Context, string) (*dto.User, error)
	IsRoot(context.Context, int64) (bool, error)
}
