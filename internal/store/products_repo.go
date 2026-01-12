package store

import "database/sql"

func ListProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query(`SELECT id, name, barcode, price, stock FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Barcode, &product.Price, &product.Stock); err != nil {
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
		`INSERT INTO products (name, barcode, price, stock) VALUES (?, ?, ?, ?)`,
		product.Name,
		product.Barcode,
		product.Price,
		product.Stock,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateProduct(db *sql.DB, product Product) error {
	_, err := db.Exec(
		`UPDATE products SET name = ?, barcode = ?, price = ?, stock = ? WHERE id = ?`,
		product.Name,
		product.Barcode,
		product.Price,
		product.Stock,
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
		`SELECT id, name, barcode, price, stock FROM products WHERE barcode = ?`,
		barcode,
	).Scan(&product.ID, &product.Name, &product.Barcode, &product.Price, &product.Stock)
	return product, err
}

func SearchProducts(db *sql.DB, query string, limit int) ([]Product, error) {
	if limit <= 0 {
		limit = 20
	}
	likeQuery := "%" + query + "%"
	rows, err := db.Query(
		`SELECT id, name, barcode, price, stock FROM products WHERE name LIKE ? OR barcode LIKE ? ORDER BY name LIMIT ?`,
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
		if err := rows.Scan(&product.ID, &product.Name, &product.Barcode, &product.Price, &product.Stock); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
