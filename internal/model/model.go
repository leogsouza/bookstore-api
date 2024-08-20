package model

type Entity interface {
	User | Book | Order
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Orders   []Order
}

type Book struct {
	ID     int
	Title  string
	Author string
	Price  float64
	Orders []Order `gorm:"many2many:order_items;"`
}

type Order struct {
	ID     int
	UserID int
	Books  []Book `gorm:"many2many:order_items;"`
	Total  float64
}

type OrderItem struct {
	OrderID  int `gorm:"primaryKey"`
	BookID   int `gorm:"primaryKey"`
	Quantity int
	Price    float64
}
