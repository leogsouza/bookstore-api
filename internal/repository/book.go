package repository

import (
	"bookstore-api/internal/model"

	"gorm.io/gorm"
)

type bookRepo[T model.Book] struct {
	db *gorm.DB
}

func NewBookRepo[T model.Book](db *gorm.DB) Reader[model.Book] {
	return &bookRepo[T]{db}
}

func (c *bookRepo[T]) Find(id int) (model.Book, error) {
	book := model.Book{}
	err := c.db.First(&book, id).Error

	return book, err
}

func (c *bookRepo[T]) FindAll() ([]model.Book, error) {
	var books = []model.Book{}
	err := c.db.Find(&books).Error
	return books, err
}

func (c *bookRepo[T]) FindByCondition(condition string, args interface{}) (*model.Book, error) {
	book := model.Book{}
	err := c.db.Where(condition, args).First(&book).Error

	return &book, err
}
