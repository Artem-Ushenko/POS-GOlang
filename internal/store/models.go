package store

type Customer struct {
	ID    int64
	Name  string
	Email string
	Phone string
}

type Product struct {
	ID            int64
	Name          string
	Barcode       string
	Quantity      int64
	PurchasePrice float64
	Price         float64
}

type SaleItem struct {
	ProductID int64
	Quantity  int64
}
