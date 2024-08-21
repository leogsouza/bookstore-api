package service

import (
	"bookstore-api/internal/model"
	"bookstore-api/internal/repository"
)

type Service[T model.Entity] interface {
	repository.Repository[T]
}

type service[T model.Entity] struct {
	repo repository.Repository[T]
}

func NewService[T model.Entity](repo repository.Repository[T]) Service[T] {
	return &service[T]{repo}
}

func (s *service[T]) Find(id int) (T, error) {
	return s.repo.Find(id)
}

func (s *service[T]) FindAll() ([]T, error) {
	return s.repo.FindAll()
}

func (s *service[T]) FindByCondition(condition string, args interface{}) (*T, error) {

	return s.repo.FindByCondition(condition, args)
}

func (s *service[T]) FindAllByCondition(condition string, args interface{}) ([]*T, error) {

	return s.repo.FindAllByCondition(condition, args)
}

func (s *service[T]) Create(entity T) (*T, error) {
	return s.repo.Create(entity)
}
