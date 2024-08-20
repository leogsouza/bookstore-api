package repository

import (
	"bookstore-api/internal/model"

	"gorm.io/gorm"
)

type userRepo[T model.User] struct {
	db *gorm.DB
}

func NewUserRepo[T model.User](db *gorm.DB) Repository[model.User] {
	return &userRepo[T]{db}
}

func (c *userRepo[T]) Find(id int) (model.User, error) {
	user := model.User{}
	err := c.db.First(&user, id).Error

	return user, err
}

func (c *userRepo[T]) FindAll() ([]model.User, error) {
	var users = []model.User{}
	err := c.db.Find(&users).Error
	return users, err
}

func (c *userRepo[T]) Create(user model.User) (int, error) {

	err := c.db.Create(&user).Error
	return user.ID, err

}
