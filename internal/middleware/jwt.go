package middleware

import (
	"bookstore-api/internal/auth"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

func JWTAuthentication(c *fiber.Ctx) error {
	token := c.Get("Bookstore-Api-Token")
	if token == "" {
		return fmt.Errorf("unauthorized")
	}

	jwtWrapper := auth.JwtWrapper{
		SecretKey: os.Getenv("JWT_SECRET"),
		Issuer:    "AuthUser",
	}

	claim, err := jwtWrapper.ValidateToken(token)

	if err != nil {
		return err
	}

	c.Locals("userEmail", claim.Email)

	return c.Next()
}
