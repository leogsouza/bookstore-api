package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
}

type Book struct {
	gorm.Model
	Title  string
	Author string
	Price  float64
	Orders []Order `gorm:"many2many:order_items;"`
}

type Order struct {
	gorm.Model
	Books []Book `gorm:"many2many:order_items;"`
	Total float64
}

type OrderItem struct {
	OrderID  int `gorm:"primaryKey"`
	BookID   int `gorm:"primaryKey"`
	Quantity int
	Price    float64
}
