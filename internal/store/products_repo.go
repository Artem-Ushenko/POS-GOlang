package store

import "database/sql"

func ListProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query(`SELECT id, name, barcode, quantity, purchase_price, price FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Barcode,
			&product.Quantity,
			&product.PurchasePrice,
			&product.Price,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func CreateProduct(db *sql.DB, product Product) (int64, error) {
	result, err := db.Exec(
		`INSERT INTO products (name, barcode, quantity, purchase_price, price) VALUES (?, ?, ?, ?, ?)`,
		product.Name,
		product.Barcode,
		product.Quantity,
		product.PurchasePrice,
		product.Price,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateProduct(db *sql.DB, product Product) error {
	_, err := db.Exec(
		`UPDATE products SET name = ?, barcode = ?, quantity = ?, purchase_price = ?, price = ? WHERE id = ?`,
		product.Name,
		product.Barcode,
		product.Quantity,
		product.PurchasePrice,
		product.Price,
		product.ID,
	)
	return err
}

func DeleteProduct(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM products WHERE id = ?`, id)
	return err
}

func GetProductByBarcode(db *sql.DB, barcode string) (Product, error) {
	var product Product
	err := db.QueryRow(
		`SELECT id, name, barcode, quantity, purchase_price, price FROM products WHERE barcode = ?`,
		barcode,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Barcode,
		&product.Quantity,
		&product.PurchasePrice,
		&product.Price,
	)
	return product, err
}

func SearchProducts(db *sql.DB, query string, limit int) ([]Product, error) {
	if limit <= 0 {
		limit = 20
	}
	likeQuery := "%" + query + "%"
	rows, err := db.Query(
		`SELECT id, name, barcode, quantity, purchase_price, price FROM products WHERE (name LIKE ? OR barcode LIKE ?) AND quantity > 0 ORDER BY name LIMIT ?`,
		likeQuery,
		likeQuery,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Barcode,
			&product.Quantity,
			&product.PurchasePrice,
			&product.Price,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
