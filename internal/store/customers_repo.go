package store

import "database/sql"

func ListCustomers(db *sql.DB) ([]Customer, error) {
	rows, err := db.Query(`SELECT id, name, email, phone FROM customers ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone); err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func CreateCustomer(db *sql.DB, customer Customer) (int64, error) {
	result, err := db.Exec(`INSERT INTO customers (name, email, phone) VALUES (?, ?, ?)`, customer.Name, customer.Email, customer.Phone)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func UpdateCustomer(db *sql.DB, customer Customer) error {
	_, err := db.Exec(`UPDATE customers SET name = ?, email = ?, phone = ? WHERE id = ?`, customer.Name, customer.Email, customer.Phone, customer.ID)
	return err
}

func DeleteCustomer(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM customers WHERE id = ?`, id)
	return err
}
