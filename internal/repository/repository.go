package repository

import "bookstore-api/internal/model"

// Repository is the interface that contains methods to interact with the database
type Repository[T model.Entity] interface {
	Reader[T]
	Writer[T]
}

type Reader[T model.Entity] interface {
	Find(id int) (T, error)
	FindAll() ([]T, error)
	FindByCondition(condition string, data interface{}) (*T, error)
	FindAllByCondition(condition string, data interface{}) ([]*T, error)
}

type Writer[T model.Entity] interface {
	Create(T) (*T, error)
}
