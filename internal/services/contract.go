package services

import (
	"automatedShop/internal/services/dto"
	"context"
)

type IShopService interface {
	// Handbook's methods
	ShowWarehousesTable(context.Context) ([]*dto.WarehousesData, error)
	CreateWarehousesItem(context.Context, *dto.WarehousesData) error
	UpdateWarehousesItem(context.Context, *dto.WarehousesData) error
	DeleteWarehousesItem(context.Context, int) error
	ShowExpenseItemsTable(context.Context) ([]*dto.ExpenseItemsData, error)
	CreateExpenseItem(context.Context, string) error
	UpdateExpenseItem(context.Context, *dto.ExpenseItemsData) error
	DeleteExpenseItem(context.Context, int) error

	// Journal's methods
	ShowChargesTable(context.Context) ([]*dto.ChargesData, error)
	CreateChargesItem(context.Context, *dto.ChargesData) error
	UpdateChargesItem(context.Context, *dto.ChargesData) error
	DeleteChargesItem(context.Context, int) error
	ShowSalesTable(context.Context) ([]*dto.SalesData, error)
	CreateSalesItem(context.Context, *dto.SalesData) error
	UpdateSalesItem(context.Context, *dto.SalesData) error
	DeleteSalesItem(context.Context, int) error

	// Report's methods
	CountMonthProfit(context.Context) (int64, error)
	GetFiveBestItems(context.Context, string, string) ([]*dto.BestItemsData, error)
}

type IAuthService interface {
	AuthoriseUser(context.Context, string, string) bool
	RegisterUser(context.Context, string, string) error
	IsRootUser(context.Context, int64) (bool, error)
}
