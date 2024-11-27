package psql

import (
	"automatedShop/internal/dataprovider"
	"automatedShop/internal/services/dto"
	"context"
	"fmt"
)

const (
	// Warehouses
	_showWarehousesTable  = `SELECT id, name, quantity, amount FROM "warehouses"`
	_insertWarehousesItem = `INSERT INTO "warehouses" (name, quantity, amount) VALUES ($1, $2, $3)`
	_updateWarehousesItem = `UPDATE "warehouses"
							 SET name = $1, quantity = $2, amount = $3
							 WHERE id = $4`
	_deleteWarehousesItem = `DELETE FROM "warehouses" WHERE id = $1`

	// Expense Items
	_showExpenseItemsTable = `SELECT id, name FROM "expense_items"`
	_insertExpenseItem     = `INSERT INTO "expense_items" (name) VALUES ($1)`
	_updateExpenseItem     = `UPDATE "expense_items"
							  SET name = $1
							  WHERE id = $2
                             `
	_deleteExpenseItem = `DELETE FROM "expense_items" WHERE id = $1`

	// Charges
	_showChargesTable  = `SELECT id, amount, charge_date, expense_item_id FROM "charges"`
	_insertChargesItem = `INSERT INTO "charges" (amount, charge_date, expense_item_id) VALUES ($1, $2, $3)`
	_updateChargesItem = `UPDATE "charges"
                              SET amount = $1, charge_date = $2, expense_item_id = $3
							  WHERE id = $4
                             `
	_deleteChargesItem = `DELETE FROM "charges" WHERE id = $1`

	// Sales
	_showSalesTable  = `SELECT id, amount, quantity, sale_date, warehouses_id FROM "sales"`
	_insertSalesItem = `INSERT INTO "sales" (amount, quantity, sale_date, warehouses_id) VALUES ($1, $2, $3, $4)`
	_updateSalesItem = `UPDATE "sales"
                              SET amount = $1, quantity = $2, sale_date = $3, warehouses_id = $4
							  WHERE id = $5
                             `
	_deleteSalesItem = `DELETE FROM "sales" WHERE id = $1`

	// Profit
	_countMonthProfit = `WITH sales_last_month AS (
							SELECT SUM(quantity * amount) AS total_sales
							FROM sales
							WHERE sale_date >= CURRENT_DATE - INTERVAL '1 month'
						),
						charges_last_month AS (
 							SELECT SUM(amount) AS total_charges
							FROM charges
							WHERE charge_date >= CURRENT_DATE - INTERVAL '1 month'
						)
						SELECT slm.total_sales - clm.total_charges AS profit
						FROM sales_last_month slm, charges_last_month clm
						`
	// 5 best items
	_showBestItems = `SELECT w.name AS warehouse_name, 
       				 	 	SUM(s.quantity * s.amount) AS total_revenue
					  	 	FROM sales s
						 	JOIN warehouses w ON s.warehouses_id = w.id
						 WHERE s.sale_date BETWEEN $1 AND $2 
						 GROUP BY w.name
						 ORDER BY total_revenue DESC
						 LIMIT 5;
                     `
)

type ShopProvider struct {
	db *dataprovider.Provider
}

func NewShopProvider(db *dataprovider.Provider) *ShopProvider {
	return &ShopProvider{db: db}
}

// Handbook's methods

func (p *ShopProvider) ShowWarehousesTable(ctx context.Context) ([]*dto.WarehousesData, error) {
	const op = "ShopRepo.ShowWarehousesTable"

	rows, err := p.db.QueryContext(ctx, _showWarehousesTable)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var warehouses []*dto.WarehousesData
	for rows.Next() {
		var warehouse dto.WarehousesData
		if err = rows.Scan(&warehouse.Id, &warehouse.Name, &warehouse.Quantity, &warehouse.Amount); err != nil {
			return nil, err
		}
		warehouses = append(warehouses, &warehouse)
	}

	return warehouses, nil
}

func (p *ShopProvider) CreateWarehousesItem(ctx context.Context, name string, quantity, amount int) error {
	const op = "ShopRepo.CreateWarehousesItem"

	_, err := p.db.ExecContext(ctx, _insertWarehousesItem, name, quantity, amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *ShopProvider) UpdateWarehousesItem(ctx context.Context, name string, quantity, amount, id int) error {
	const op = "ShopRepo.UpdateWarehousesItem"

	_, err := p.db.ExecContext(ctx, _updateWarehousesItem, name, quantity, amount, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *ShopProvider) DeleteWarehousesItem(ctx context.Context, id int) error {
	const op = "ShopRepo.DeleteWarehousesItem"

	_, err := p.db.ExecContext(ctx, _deleteWarehousesItem, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *ShopProvider) ShowExpenseItemsTable(ctx context.Context) ([]*dto.ExpenseItemsData, error) {
	const op = "ShopRepo.ShowExpenseItemsTable"

	rows, err := p.db.QueryContext(ctx, _showExpenseItemsTable)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var exItems []*dto.ExpenseItemsData
	for rows.Next() {
		var exItem dto.ExpenseItemsData
		if err = rows.Scan(&exItem.Id, &exItem.Name); err != nil {
			return nil, err
		}
		exItems = append(exItems, &exItem)
	}

	return exItems, nil
}

func (p *ShopProvider) CreateExpenseItem(ctx context.Context, name string) error {
	const op = "ShopRepo.CreateExpenseItem"

	_, err := p.db.ExecContext(ctx, _insertExpenseItem, name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *ShopProvider) UpdateExpenseItem(ctx context.Context, name string, id int) error {
	const op = "ShopRepo.UpdateExpenseItem"

	_, err := p.db.ExecContext(ctx, _updateExpenseItem, name, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *ShopProvider) DeleteExpenseItem(ctx context.Context, id int) error {
	const op = "ShopRepo.DeleteExpenseItem"

	_, err := p.db.ExecContext(ctx, _deleteExpenseItem, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Journal's methods
func (p *ShopProvider) ShowChargesTable(ctx context.Context) ([]*dto.ChargesData, error) {
	const op = "ShopRepo.ShowChargesTable"

	rows, err := p.db.QueryContext(ctx, _showChargesTable)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var chargesItems []*dto.ChargesData
	for rows.Next() {
		var chargesItem dto.ChargesData
		if err = rows.Scan(&chargesItem.Id, &chargesItem.Amount, &chargesItem.ChargeDate,
			&chargesItem.ExpenseItemId); err != nil {
			return nil, err
		}
		chargesItems = append(chargesItems, &chargesItem)
	}

	return chargesItems, nil
}

func (p *ShopProvider) CreateChargesItem(ctx context.Context, data *dto.ChargesData) error {
	const op = "ShopRepo.CreateChargesItem"

	_, err := p.db.ExecContext(ctx, _insertChargesItem, data.Amount, data.ChargeDate, data.ExpenseItemId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *ShopProvider) UpdateChargesItem(ctx context.Context, data *dto.ChargesData) error {
	const op = "ShopRepo.UpdateChargesItem"

	_, err := p.db.ExecContext(ctx, _updateChargesItem, data.Amount, data.ChargeDate, data.ExpenseItemId, data.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *ShopProvider) DeleteChargesItem(ctx context.Context, id int) error {
	const op = "ShopRepo.DeleteChargesItem"

	_, err := p.db.ExecContext(ctx, _deleteChargesItem, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *ShopProvider) ShowSalesTable(ctx context.Context) ([]*dto.SalesData, error) {
	const op = "ShopRepo.ShowSalesTable"

	rows, err := p.db.QueryContext(ctx, _showSalesTable)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var salesItems []*dto.SalesData
	for rows.Next() {
		var salesItem dto.SalesData
		if err = rows.Scan(&salesItem.Id, &salesItem.Amount, &salesItem.Quantity, &salesItem.SaleDate,
			&salesItem.WarehousesId); err != nil {
			return nil, err
		}
		salesItems = append(salesItems, &salesItem)
	}

	return salesItems, nil
}

func (p *ShopProvider) CreateSalesItem(ctx context.Context, data *dto.SalesData) error {
	const op = "ShopRepo.CreateSalesItem"

	_, err := p.db.ExecContext(ctx, _insertSalesItem, data.Amount, data.Quantity, data.SaleDate, data.WarehousesId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *ShopProvider) UpdateSalesItem(ctx context.Context, data *dto.SalesData) error {
	const op = "ShopRepo.UpdateSalesItem"

	_, err := p.db.ExecContext(ctx, _updateSalesItem, data.Amount, data.Quantity, data.SaleDate, data.WarehousesId, data.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *ShopProvider) DeleteSalesItem(ctx context.Context, id int) error {
	const op = "ShopRepo.DeleteSalesItem"

	_, err := p.db.ExecContext(ctx, _deleteSalesItem, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Report's methods

func (p *ShopProvider) CountMonthProfit(ctx context.Context) (int64, error) {
	const op = "ShopRepo.CountMonthProfit"

	var profit int64

	err := p.db.GetContext(ctx, &profit, _countMonthProfit)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return profit, nil
}

func (p *ShopProvider) GetFiveBestItems(ctx context.Context, from string, to string) ([]*dto.BestItemsData, error) {
	const op = "ShopRepo.GetFiveBestItems"

	rows, err := p.db.QueryContext(ctx, _showBestItems, from, to)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var items []*dto.BestItemsData
	for rows.Next() {
		var item dto.BestItemsData
		if err = rows.Scan(&item.Name, &item.TotalRevenue); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}
