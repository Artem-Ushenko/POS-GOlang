package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
			if err := saveStore(dataFile, store); err != nil {
				fmt.Printf("Failed to save data: %v\n", err)
			}
		case 3:
			editCustomer(reader, store)
			if err := saveStore(dataFile, store); err != nil {
				fmt.Printf("Failed to save data: %v\n", err)
			}
		case 4:
			removeCustomer(reader, store)
			if err := saveStore(dataFile, store); err != nil {
				fmt.Printf("Failed to save data: %v\n", err)
			}
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
			if err := saveStore(dataFile, store); err != nil {
				fmt.Printf("Failed to save data: %v\n", err)
			}
		case 3:
			editProduct(reader, store)
			if err := saveStore(dataFile, store); err != nil {
				fmt.Printf("Failed to save data: %v\n", err)
			}
		case 4:
			removeProduct(reader, store)
			if err := saveStore(dataFile, store); err != nil {
				fmt.Printf("Failed to save data: %v\n", err)
			}
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
			if err := saveStore(dataFile, store); err != nil {
				fmt.Printf("Failed to save data: %v\n", err)
			}
		case 3:
			editSale(reader, store)
			if err := saveStore(dataFile, store); err != nil {
				fmt.Printf("Failed to save data: %v\n", err)
			}
		case 4:
			removeSale(reader, store)
			if err := saveStore(dataFile, store); err != nil {
				fmt.Printf("Failed to save data: %v\n", err)
			}
		case 0:
			return
		default:
			fmt.Println("Invalid option.")
		}
	}
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
