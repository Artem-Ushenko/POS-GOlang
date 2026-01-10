package main

import "time"

const dataFile = "data.json"

type Customer struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

type Sale struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	ProductID  int       `json:"product_id"`
	Quantity   int       `json:"quantity"`
	Total      float64   `json:"total"`
	CreatedAt  time.Time `json:"created_at"`
}

type Store struct {
	NextCustomerID int        `json:"next_customer_id"`
	NextProductID  int        `json:"next_product_id"`
	NextSaleID     int        `json:"next_sale_id"`
	Customers      []Customer `json:"customers"`
	Products       []Product  `json:"products"`
	Sales          []Sale     `json:"sales"`
}
