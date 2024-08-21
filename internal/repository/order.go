package repository

import (
	"bookstore-api/internal/model"

	"gorm.io/gorm"
)

type orderRepo[T model.Order] struct {
	db *gorm.DB
}

func NewOrderRepo[T model.Order](db *gorm.DB) Repository[model.Order] {
	return &orderRepo[T]{db}
}

func (c *orderRepo[T]) Find(id int) (model.Order, error) {
	order := model.Order{}
	err := c.db.Preload("OrderItems").First(&order, id).Error

	return order, err
}

func (c *orderRepo[T]) FindAll() ([]model.Order, error) {
	var orders = []model.Order{}
	err := c.db.Preload("OrderItems").Find(&orders).Error
	return orders, err
}

func (c *orderRepo[T]) FindByCondition(condition string, args interface{}) (*model.Order, error) {
	order := model.Order{}
	err := c.db.Where(condition, args).Preload("OrderItems").First(&order).Error

	return &order, err
}

func (c *orderRepo[T]) FindAllByCondition(condition string, args interface{}) ([]*model.Order, error) {
	orders := []*model.Order{}
	err := c.db.Where(condition, args).Preload("OrderItems").Find(&orders).Error

	return orders, err
}

func (c *orderRepo[T]) Create(order model.Order) (*model.Order, error) {

	result := c.db.Create(&order)

	return &order, result.Error

}
