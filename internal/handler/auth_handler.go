package handler

import (
	"bookstore-api/internal/auth"
	"bookstore-api/internal/model"
	"bookstore-api/internal/service"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type authHandler struct {
	service service.Service[model.User]
}

func NewAuthHandler(service service.Service[model.User]) *authHandler {
	return &authHandler{service}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}

type genericResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResponse{
		Type:    "error",
		Message: "invalid credentials",
	})
}

func (h *authHandler) Authenticate(c *fiber.Ctx) error {
	var authParams AuthParams

	if err := c.BodyParser(&authParams); err != nil {
		return err
	}

	user, err := h.service.FindByCondition("email = ?", authParams.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return invalidCredentials(c)
		}
		return err
	}

	if !IsValidPassword(user.Password, authParams.Password) {
		return invalidCredentials(c)
	}

	jwtWrapper := auth.JwtWrapper{
		SecretKey:       os.Getenv("JWT_SECRET"),
		Issuer:          "AuthUser",
		ExpirationHours: 4,
	}

	tokenStr, err := jwtWrapper.GenerateToken(user.Email)
	if err != nil {
		return fmt.Errorf("cannot authenticate")
	}

	resp := &authResponse{
		User:  user,
		Token: tokenStr,
	}
	return c.JSON(resp)
}
