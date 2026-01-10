package store

import "database/sql"

func ListProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query(`SELECT id, name, price, stock FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock); err != nil {
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
	result, err := db.Exec(`INSERT INTO products (name, price, stock) VALUES (?, ?, ?)`, product.Name, product.Price, product.Stock)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateProduct(db *sql.DB, product Product) error {
	_, err := db.Exec(`UPDATE products SET name = ?, price = ?, stock = ? WHERE id = ?`, product.Name, product.Price, product.Stock, product.ID)
	return err
}

func DeleteProduct(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM products WHERE id = ?`, id)
	return err
}
