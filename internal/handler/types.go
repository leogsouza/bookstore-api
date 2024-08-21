package handler

import (
	"bookstore-api/internal/model"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12
)

type userRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8"`
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

type orderRequest struct {
	Items []orderItemRequest `json:"items" validate:"required,gte=1"`
}

type orderItemRequest struct {
	BookID   int     `json:"book_id" validate:"required"`
	Quantity int     `json:"quantity" validate:"required,gte=1"`
	Price    float64 `json:"price" validate:"required,gt=0"`
}

func NewOrderFromRequest(req orderRequest) *model.Order {

	books := []model.OrderItem{}

	var total float64
	for _, item := range req.Items {
		book := model.OrderItem{BookID: item.BookID, Price: item.Price, Quantity: item.Quantity}

		books = append(books, book)
		total += (item.Price * float64(item.Quantity))
	}

	return &model.Order{
		OrderItems: books,
		Total:      total,
	}
}
