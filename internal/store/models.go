package store

type Customer struct {
	ID    int64
	Name  string
	Email string
	Phone string
}

type Product struct {
	ID      int64
	Name    string
	Barcode string
	Price   float64
	Stock   int64
}

type SaleItem struct {
	ProductID int64
	Quantity  int64
}
