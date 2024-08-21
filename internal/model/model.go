package model

type Entity interface {
	User | Book | Order
}

type User struct {
	ID       int
	Name     string `gorm:"not null;"`
	Email    string `gorm:"uniqueIndex; not null"`
	Password string `gorm:"not null" json:"-"`
	Orders   []Order
}

type Book struct {
	ID     int
	Title  string `gorm:"not null"`
	Author string
	Price  float64 `gorm:"not null"`
	Orders []Order `gorm:"many2many:order_items;"`
}

type Order struct {
	ID         int
	UserID     int     `gorm:"not null"`
	Total      float64 `gorm:"not null"`
	OrderItems []OrderItem
}

type OrderItem struct {
	OrderID  int     `gorm:"primaryKey"`
	BookID   int     `gorm:"primaryKey"`
	Quantity int     `gorm:"not null"`
	Price    float64 `gorm:"not null"`
}
