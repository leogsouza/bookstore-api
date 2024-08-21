package handler

import (
	"bookstore-api/internal/model"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12
)

type userRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserFromRequest(req userRequest) (*model.User, error) {
	encpw, err := req.HashPassword()
	if err != nil {
		return nil, err
	}

	return &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(encpw),
	}, nil
}

func (u *userRequest) HashPassword() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(u.Password), bcryptCost)
}

func IsValidPassword(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}
