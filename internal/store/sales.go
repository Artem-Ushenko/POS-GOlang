package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrInsufficientStock = errors.New("insufficient stock")

func CreateSale(db *sql.DB, customerID *int64, items []SaleItem) (int64, error) {
	if customerID == nil {
		return 0, errors.New("customer id is required")
	}
	if len(items) == 0 {
		return 0, errors.New("sale requires at least one item")
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	itemPrices := make(map[int64]float64, len(items))
	for _, item := range items {
		if item.Quantity <= 0 {
			return 0, errors.New("quantity must be greater than zero")
		}

		var price float64
		var quantity int64
		err := tx.QueryRow(
			`SELECT price, quantity FROM products WHERE id = ?`,
			item.ProductID,
		).Scan(&price, &quantity)
		if err != nil {
			return 0, err
		}
		if quantity < item.Quantity {
			return 0, fmt.Errorf("%w for product %d", ErrInsufficientStock, item.ProductID)
		}

		result, err := tx.Exec(
			`UPDATE products SET quantity = quantity - ? WHERE id = ? AND quantity >= ?`,
			item.Quantity,
			item.ProductID,
			item.Quantity,
		)
		if err != nil {
			return 0, err
		}
		updated, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		if updated == 0 {
			return 0, fmt.Errorf("%w for product %d", ErrInsufficientStock, item.ProductID)
		}

		itemPrices[item.ProductID] = price
	}

	result, err := tx.Exec(
		`INSERT INTO sales (customer_id, created_at) VALUES (?, ?)`,
		*customerID,
		time.Now(),
	)
	if err != nil {
		return 0, err
	}
	saleID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, item := range items {
		price := itemPrices[item.ProductID]
		_, err := tx.Exec(
			`INSERT INTO sale_items (sale_id, product_id, quantity, price) VALUES (?, ?, ?, ?)`,
			saleID,
			item.ProductID,
			item.Quantity,
			price,
		)
		if err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return saleID, nil
}
