package services

import (
	"automatedShop/internal/repository"
	"automatedShop/internal/services/dto"
	"context"
	"fmt"
	"log/slog"
)

type ShopService struct {
	l        *slog.Logger
	ShopRepo repository.IShopRepository
}

func NewShopService(repo repository.IShopRepository) *ShopService {
	var l *slog.Logger

	return &ShopService{
		l:        l,
		ShopRepo: repo,
	}
}

func (s *ShopService) ShowWarehousesTable(ctx context.Context) ([]*dto.WarehousesData, error) {
	const op = "ShopService.ShowWarehousesTable"

	res, err := s.ShopRepo.ShowWarehousesTable(ctx)
	if err != nil {
		return nil, fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	return res, nil
}

func (s *ShopService) CreateWarehousesItem(ctx context.Context, data *dto.WarehousesData) error {
	const op = "ShopService.CreateWarehousesItem"

	err := s.ShopRepo.CreateWarehousesItem(ctx, data.Name, data.Quantity, data.Amount)
	if err != nil {
		return fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	fmt.Printf("%v: warehouses item inserted successfully", op)
	return nil
}

func (s *ShopService) UpdateWarehousesItem(ctx context.Context, data *dto.WarehousesData) error {
	const op = "ShopService.UpdateWarehousesItem"

	err := s.ShopRepo.UpdateWarehousesItem(ctx, data.Name, data.Quantity, data.Amount, data.Id)
	if err != nil {
		return fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	fmt.Printf("%v: warehouses item updated successfully", op)
	return nil
}

func (s *ShopService) DeleteWarehousesItem(ctx context.Context, id int) error {
	const op = "ShopService.DeleteWarehousesItem"

	err := s.ShopRepo.DeleteWarehousesItem(ctx, id)
	if err != nil {
		return fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	fmt.Printf("%v: warehouses item deleted successfully", op)
	return nil
}

func (s *ShopService) ShowExpenseItemsTable(ctx context.Context) ([]*dto.ExpenseItemsData, error) {
	const op = "ShopService.ShowExpenseItemsTable"

	res, err := s.ShopRepo.ShowExpenseItemsTable(ctx)
	if err != nil {
		return nil, fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	return res, nil
}

func (s *ShopService) CreateExpenseItem(ctx context.Context, name string) error {
	const op = "ShopService.CreateExpenseItem"

	err := s.ShopRepo.CreateExpenseItem(ctx, name)
	if err != nil {
		return fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	fmt.Printf("%v: expense item inserted successfully", op)
	return nil
}

func (s *ShopService) UpdateExpenseItem(ctx context.Context, data *dto.ExpenseItemsData) error {
	const op = "ShopService.UpdateExpenseItem"

	err := s.ShopRepo.UpdateExpenseItem(ctx, data.Name, data.Id)
	if err != nil {
		return fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	fmt.Printf("%v: expense item updated successfully", op)
	return nil
}

func (s *ShopService) DeleteExpenseItem(ctx context.Context, id int) error {
	const op = "ShopService.DeleteExpenseItem"

	err := s.ShopRepo.DeleteExpenseItem(ctx, id)
	if err != nil {
		return fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	fmt.Printf("%v: expense item deleted successfully", op)
	return nil
}

// Journal's methods
func (s *ShopService) ShowChargesTable(ctx context.Context) ([]*dto.ChargesData, error) {
	const op = "ShopService.ShowChargesTable"

	res, err := s.ShopRepo.ShowChargesTable(ctx)
	if err != nil {
		return nil, fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	return res, nil
}

func (s *ShopService) CreateChargesItem(ctx context.Context, data *dto.ChargesData) error {
	const op = "ShopService.CreateChargesItem"

	err := s.ShopRepo.CreateChargesItem(ctx, data)
	if err != nil {
		return fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	return nil
}

func (s *ShopService) UpdateChargesItem(ctx context.Context, data *dto.ChargesData) error {
	const op = "ShopService.UpdateChargesItem"

	err := s.ShopRepo.UpdateChargesItem(ctx, data)
	if err != nil {
		return fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	return nil
}

func (s *ShopService) DeleteChargesItem(ctx context.Context, id int) error {
	const op = "ShopService.DeleteChargesItem"

	err := s.ShopRepo.DeleteChargesItem(ctx, id)
	if err != nil {
		return fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	return nil
}

func (s *ShopService) ShowSalesTable(ctx context.Context) ([]*dto.SalesData, error) {
	const op = "ShopService.ShowSalesTable"

	res, err := s.ShopRepo.ShowSalesTable(ctx)
	if err != nil {
		return nil, fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	return res, nil
}

func (s *ShopService) CreateSalesItem(ctx context.Context, data *dto.SalesData) error {
	const op = "ShopService.CreateSalesItem"

	err := s.ShopRepo.CreateSalesItem(ctx, data)
	if err != nil {
		return fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	return nil
}

func (s *ShopService) UpdateSalesItem(ctx context.Context, data *dto.SalesData) error {
	const op = "ShopService.UpdateSalesItem"

	err := s.ShopRepo.UpdateSalesItem(ctx, data)
	if err != nil {
		return fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	return nil
}

func (s *ShopService) DeleteSalesItem(ctx context.Context, id int) error {
	const op = "ShopService.DeleteSalesItem"

	err := s.ShopRepo.DeleteSalesItem(ctx, id)
	if err != nil {
		return fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	return nil
}

// Report's methods

func (s *ShopService) CountMonthProfit(ctx context.Context) (int64, error) {
	const op = "ShopService.CountMonthProfit"

	profit, err := s.ShopRepo.CountMonthProfit(ctx)
	if err != nil {
		return -1, fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	fmt.Printf("%v: profit counted successfully", op)
	return profit, nil
}

func (s *ShopService) GetFiveBestItems(ctx context.Context, from string, to string) ([]*dto.BestItemsData, error) {
	const op = "ShopService.GetFiveBestItems"

	res, err := s.ShopRepo.GetFiveBestItems(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("error occurred in: %v: %v", op, err)
	}

	return res, nil
}
