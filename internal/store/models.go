package store

type Customer struct {
	ID    int64
	Name  string
	Email string
	Phone string
}

type Product struct {
	ID    int64
	Name  string
	Price float64
	Stock int64
}
