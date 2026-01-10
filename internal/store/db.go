package store

import "database/sql"

func OpenDB(path string) (*sql.DB, error) {
	return sql.Open("sqlite", path)
}

func Migrate(db *sql.DB) error {
	if _, err := db.Exec(`PRAGMA journal_mode = DELETE;`); err != nil {
		return err
	}

	statements := []string{
		`CREATE TABLE IF NOT EXISTS customers (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			phone TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			price REAL NOT NULL,
			stock INTEGER NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS sales (
			id INTEGER PRIMARY KEY,
			customer_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL,
			FOREIGN KEY(customer_id) REFERENCES customers(id)
		);`,
		`CREATE TABLE IF NOT EXISTS sale_items (
			id INTEGER PRIMARY KEY,
			sale_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			price REAL NOT NULL,
			FOREIGN KEY(sale_id) REFERENCES sales(id),
			FOREIGN KEY(product_id) REFERENCES products(id)
		);`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}
