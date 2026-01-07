package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

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

func main() {
	store, err := loadStore(dataFile)
	if err != nil {
		fmt.Printf("Failed to load data: %v\n", err)
		store = newStore()
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n=== POS System ===")
		fmt.Println("1) Customers")
		fmt.Println("2) Products")
		fmt.Println("3) Sales")
		fmt.Println("0) Save & Exit")

		switch readInt(reader, "Select option: ") {
		case 1:
			manageCustomers(reader, store)
		case 2:
			manageProducts(reader, store)
		case 3:
			manageSales(reader, store)
		case 0:
			if err := saveStore(dataFile, store); err != nil {
				fmt.Printf("Failed to save data: %v\n", err)
			} else {
				fmt.Println("Data saved.")
			}
			return
		default:
			fmt.Println("Invalid option.")
		}
	}
}

func newStore() *Store {
	return &Store{
		NextCustomerID: 1,
		NextProductID:  1,
		NextSaleID:     1,
	}
}

func loadStore(path string) (*Store, error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return newStore(), nil
		}
		return nil, err
	}
	defer file.Close()

	var store Store
	if err := json.NewDecoder(file).Decode(&store); err != nil {
		return nil, err
	}

	return &store, nil
}

func saveStore(path string, store *Store) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(store)
}

func manageCustomers(reader *bufio.Reader, store *Store) {
	for {
		fmt.Println("\n--- Customers ---")
		fmt.Println("1) List")
		fmt.Println("2) Add")
		fmt.Println("3) Edit")
		fmt.Println("4) Remove")
		fmt.Println("0) Back")

		switch readInt(reader, "Select option: ") {
		case 1:
			listCustomers(store)
		case 2:
			addCustomer(reader, store)
			saveStore(dataFile, store)
		case 3:
			editCustomer(reader, store)
			saveStore(dataFile, store)
		case 4:
			removeCustomer(reader, store)
			saveStore(dataFile, store)
		case 0:
			return
		default:
			fmt.Println("Invalid option.")
		}
	}
}

func manageProducts(reader *bufio.Reader, store *Store) {
	for {
		fmt.Println("\n--- Products ---")
		fmt.Println("1) List")
		fmt.Println("2) Add")
		fmt.Println("3) Edit")
		fmt.Println("4) Remove")
		fmt.Println("0) Back")

		switch readInt(reader, "Select option: ") {
		case 1:
			listProducts(store)
		case 2:
			addProduct(reader, store)
			saveStore(dataFile, store)
		case 3:
			editProduct(reader, store)
			saveStore(dataFile, store)
		case 4:
			removeProduct(reader, store)
			saveStore(dataFile, store)
		case 0:
			return
		default:
			fmt.Println("Invalid option.")
		}
	}
}

func manageSales(reader *bufio.Reader, store *Store) {
	for {
		fmt.Println("\n--- Sales ---")
		fmt.Println("1) List")
		fmt.Println("2) Add")
		fmt.Println("3) Edit")
		fmt.Println("4) Remove")
		fmt.Println("0) Back")

		switch readInt(reader, "Select option: ") {
		case 1:
			listSales(store)
		case 2:
			addSale(reader, store)
			saveStore(dataFile, store)
		case 3:
			editSale(reader, store)
			saveStore(dataFile, store)
		case 4:
			removeSale(reader, store)
			saveStore(dataFile, store)
		case 0:
			return
		default:
			fmt.Println("Invalid option.")
		}
	}
}

func listCustomers(store *Store) {
	fmt.Println("\nCustomers:")
	if len(store.Customers) == 0 {
		fmt.Println("No customers found.")
		return
	}
	for _, customer := range store.Customers {
		fmt.Printf("ID: %d | %s | %s | %s\n", customer.ID, customer.Name, customer.Email, customer.Phone)
	}
}

func addCustomer(reader *bufio.Reader, store *Store) {
	name := readString(reader, "Name: ")
	email := readString(reader, "Email: ")
	phone := readString(reader, "Phone: ")

	store.Customers = append(store.Customers, Customer{
		ID:    store.NextCustomerID,
		Name:  name,
		Email: email,
		Phone: phone,
	})
	store.NextCustomerID++

	fmt.Println("Customer added.")
}

func editCustomer(reader *bufio.Reader, store *Store) {
	id := readInt(reader, "Customer ID to edit: ")
	index := findCustomerIndex(store, id)
	if index == -1 {
		fmt.Println("Customer not found.")
		return
	}

	customer := store.Customers[index]
	name := readOptionalString(reader, fmt.Sprintf("Name [%s]: ", customer.Name))
	email := readOptionalString(reader, fmt.Sprintf("Email [%s]: ", customer.Email))
	phone := readOptionalString(reader, fmt.Sprintf("Phone [%s]: ", customer.Phone))

	if name != "" {
		customer.Name = name
	}
	if email != "" {
		customer.Email = email
	}
	if phone != "" {
		customer.Phone = phone
	}

	store.Customers[index] = customer
	fmt.Println("Customer updated.")
}

func removeCustomer(reader *bufio.Reader, store *Store) {
	id := readInt(reader, "Customer ID to remove: ")
	index := findCustomerIndex(store, id)
	if index == -1 {
		fmt.Println("Customer not found.")
		return
	}

	if hasCustomerSales(store, id) {
		fmt.Println("Cannot remove customer with existing sales.")
		return
	}

	store.Customers = append(store.Customers[:index], store.Customers[index+1:]...)
	fmt.Println("Customer removed.")
}

func listProducts(store *Store) {
	fmt.Println("\nProducts:")
	if len(store.Products) == 0 {
		fmt.Println("No products found.")
		return
	}
	for _, product := range store.Products {
		fmt.Printf("ID: %d | %s | $%.2f | Stock: %d\n", product.ID, product.Name, product.Price, product.Stock)
	}
}

func addProduct(reader *bufio.Reader, store *Store) {
	name := readString(reader, "Name: ")
	price := readFloat(reader, "Price: ")
	stock := readInt(reader, "Stock: ")

	store.Products = append(store.Products, Product{
		ID:    store.NextProductID,
		Name:  name,
		Price: price,
		Stock: stock,
	})
	store.NextProductID++

	fmt.Println("Product added.")
}

func editProduct(reader *bufio.Reader, store *Store) {
	id := readInt(reader, "Product ID to edit: ")
	index := findProductIndex(store, id)
	if index == -1 {
		fmt.Println("Product not found.")
		return
	}

	product := store.Products[index]
	name := readOptionalString(reader, fmt.Sprintf("Name [%s]: ", product.Name))
	priceStr := readOptionalString(reader, fmt.Sprintf("Price [%.2f]: ", product.Price))
	stockStr := readOptionalString(reader, fmt.Sprintf("Stock [%d]: ", product.Stock))

	if name != "" {
		product.Name = name
	}
	if priceStr != "" {
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			fmt.Println("Invalid price; keeping existing.")
		} else {
			product.Price = price
		}
	}
	if stockStr != "" {
		stock, err := strconv.Atoi(stockStr)
		if err != nil {
			fmt.Println("Invalid stock; keeping existing.")
		} else {
			product.Stock = stock
		}
	}

	store.Products[index] = product
	fmt.Println("Product updated.")
}

func removeProduct(reader *bufio.Reader, store *Store) {
	id := readInt(reader, "Product ID to remove: ")
	index := findProductIndex(store, id)
	if index == -1 {
		fmt.Println("Product not found.")
		return
	}

	if hasProductSales(store, id) {
		fmt.Println("Cannot remove product with existing sales.")
		return
	}

	store.Products = append(store.Products[:index], store.Products[index+1:]...)
	fmt.Println("Product removed.")
}

func listSales(store *Store) {
	fmt.Println("\nSales:")
	if len(store.Sales) == 0 {
		fmt.Println("No sales found.")
		return
	}
	for _, sale := range store.Sales {
		customerName := lookupCustomerName(store, sale.CustomerID)
		productName := lookupProductName(store, sale.ProductID)
		fmt.Printf("ID: %d | Customer: %s | Product: %s | Qty: %d | Total: $%.2f | Date: %s\n",
			sale.ID,
			customerName,
			productName,
			sale.Quantity,
			sale.Total,
			sale.CreatedAt.Format(time.RFC3339),
		)
	}
}

func addSale(reader *bufio.Reader, store *Store) {
	if len(store.Customers) == 0 || len(store.Products) == 0 {
		fmt.Println("Please add at least one customer and product before creating sales.")
		return
	}

	listCustomers(store)
	customerID := readInt(reader, "Customer ID: ")
	if findCustomerIndex(store, customerID) == -1 {
		fmt.Println("Customer not found.")
		return
	}

	listProducts(store)
	productID := readInt(reader, "Product ID: ")
	productIndex := findProductIndex(store, productID)
	if productIndex == -1 {
		fmt.Println("Product not found.")
		return
	}

	quantity := readInt(reader, "Quantity: ")
	if quantity <= 0 {
		fmt.Println("Quantity must be greater than zero.")
		return
	}

	product := store.Products[productIndex]
	if product.Stock < quantity {
		fmt.Println("Not enough stock.")
		return
	}

	product.Stock -= quantity
	store.Products[productIndex] = product

	total := product.Price * float64(quantity)
	store.Sales = append(store.Sales, Sale{
		ID:         store.NextSaleID,
		CustomerID: customerID,
		ProductID:  productID,
		Quantity:   quantity,
		Total:      total,
		CreatedAt:  time.Now(),
	})
	store.NextSaleID++

	fmt.Println("Sale added.")
}

func editSale(reader *bufio.Reader, store *Store) {
	id := readInt(reader, "Sale ID to edit: ")
	index := findSaleIndex(store, id)
	if index == -1 {
		fmt.Println("Sale not found.")
		return
	}

	sale := store.Sales[index]
	oldProductIndex := findProductIndex(store, sale.ProductID)
	if oldProductIndex != -1 {
		oldProduct := store.Products[oldProductIndex]
		oldProduct.Stock += sale.Quantity
		store.Products[oldProductIndex] = oldProduct
	}

	listCustomers(store)
	customerID := readInt(reader, fmt.Sprintf("Customer ID [%d]: ", sale.CustomerID))
	if customerID == 0 {
		customerID = sale.CustomerID
	}
	if findCustomerIndex(store, customerID) == -1 {
		fmt.Println("Customer not found.")
		restoreSaleStock(store, sale)
		return
	}

	listProducts(store)
	productID := readInt(reader, fmt.Sprintf("Product ID [%d]: ", sale.ProductID))
	if productID == 0 {
		productID = sale.ProductID
	}
	productIndex := findProductIndex(store, productID)
	if productIndex == -1 {
		fmt.Println("Product not found.")
		restoreSaleStock(store, sale)
		return
	}

	quantity := readInt(reader, fmt.Sprintf("Quantity [%d]: ", sale.Quantity))
	if quantity == 0 {
		quantity = sale.Quantity
	}
	if quantity <= 0 {
		fmt.Println("Quantity must be greater than zero.")
		restoreSaleStock(store, sale)
		return
	}

	product := store.Products[productIndex]
	if product.Stock < quantity {
		fmt.Println("Not enough stock.")
		restoreSaleStock(store, sale)
		return
	}

	product.Stock -= quantity
	store.Products[productIndex] = product

	sale.CustomerID = customerID
	sale.ProductID = productID
	sale.Quantity = quantity
	sale.Total = product.Price * float64(quantity)
	store.Sales[index] = sale

	fmt.Println("Sale updated.")
}

func restoreSaleStock(store *Store, sale Sale) {
	productIndex := findProductIndex(store, sale.ProductID)
	if productIndex == -1 {
		return
	}
	product := store.Products[productIndex]
	product.Stock -= sale.Quantity
	store.Products[productIndex] = product
}

func removeSale(reader *bufio.Reader, store *Store) {
	id := readInt(reader, "Sale ID to remove: ")
	index := findSaleIndex(store, id)
	if index == -1 {
		fmt.Println("Sale not found.")
		return
	}

	sale := store.Sales[index]
	productIndex := findProductIndex(store, sale.ProductID)
	if productIndex != -1 {
		product := store.Products[productIndex]
		product.Stock += sale.Quantity
		store.Products[productIndex] = product
	}

	store.Sales = append(store.Sales[:index], store.Sales[index+1:]...)
	fmt.Println("Sale removed.")
}

func findCustomerIndex(store *Store, id int) int {
	for i, customer := range store.Customers {
		if customer.ID == id {
			return i
		}
	}
	return -1
}

func findProductIndex(store *Store, id int) int {
	for i, product := range store.Products {
		if product.ID == id {
			return i
		}
	}
	return -1
}

func findSaleIndex(store *Store, id int) int {
	for i, sale := range store.Sales {
		if sale.ID == id {
			return i
		}
	}
	return -1
}

func hasCustomerSales(store *Store, customerID int) bool {
	for _, sale := range store.Sales {
		if sale.CustomerID == customerID {
			return true
		}
	}
	return false
}

func hasProductSales(store *Store, productID int) bool {
	for _, sale := range store.Sales {
		if sale.ProductID == productID {
			return true
		}
	}
	return false
}

func lookupCustomerName(store *Store, id int) string {
	for _, customer := range store.Customers {
		if customer.ID == id {
			return customer.Name
		}
	}
	return fmt.Sprintf("Unknown (%d)", id)
}

func lookupProductName(store *Store, id int) string {
	for _, product := range store.Products {
		if product.ID == id {
			return product.Name
		}
	}
	return fmt.Sprintf("Unknown (%d)", id)
}

func readString(reader *bufio.Reader, prompt string) string {
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Input error. Please try again.")
			continue
		}
		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println("Value cannot be empty.")
			continue
		}
		return input
	}
}

func readOptionalString(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(input)
}

func readInt(reader *bufio.Reader, prompt string) int {
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Input error. Please try again.")
			continue
		}
		input = strings.TrimSpace(input)
		if input == "" {
			return 0
		}
		value, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid number. Please try again.")
			continue
		}
		return value
	}
}

func readFloat(reader *bufio.Reader, prompt string) float64 {
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Input error. Please try again.")
			continue
		}
		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println("Value cannot be empty.")
			continue
		}
		value, err := strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Println("Invalid number. Please try again.")
			continue
		}
		return value
	}
}
