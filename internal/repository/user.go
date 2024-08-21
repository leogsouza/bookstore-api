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
	err := c.db.Preload("Orders").First(&user, id).Error

	return user, err
}

func (c *userRepo[T]) FindAll() ([]model.User, error) {
	var users = []model.User{}
	err := c.db.Preload("Orders").Find(&users).Error
	return users, err
}

func (c *userRepo[T]) FindByCondition(condition string, args interface{}) (*model.User, error) {
	user := model.User{}
	err := c.db.Where(condition, args).Preload("Orders").First(&user).Error

	return &user, err
}

func (c *userRepo[T]) FindAllByCondition(condition string, args interface{}) ([]*model.User, error) {
	users := []*model.User{}
	err := c.db.Where(condition, args).Preload("Orders").Find(&users).Error

	return users, err
}

func (c *userRepo[T]) Create(user model.User) (*model.User, error) {

	result := c.db.Create(&user)

	return &user, result.Error

}
