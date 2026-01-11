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
			barcode TEXT,
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
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_products_barcode ON products(barcode);`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}

	hasBarcode, err := columnExists(db, "products", "barcode")
	if err != nil {
		return err
	}
	if !hasBarcode {
		if _, err := db.Exec(`ALTER TABLE products ADD COLUMN barcode TEXT`); err != nil {
			return err
		}
	}

	return nil
}

func columnExists(db *sql.DB, tableName, columnName string) (bool, error) {
	rows, err := db.Query(`PRAGMA table_info(` + tableName + `)`)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cid       int
			name      string
			colType   string
			notNull   int
			defaultOK sql.NullString
			pk        int
		)
		if err := rows.Scan(&cid, &name, &colType, &notNull, &defaultOK, &pk); err != nil {
			return false, err
		}
		if name == columnName {
			return true, nil
		}
	}
	if err := rows.Err(); err != nil {
		return false, err
	}
	return false, nil
}
