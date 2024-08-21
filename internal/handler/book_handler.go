package handler

import (
	"bookstore-api/internal/model"
	"bookstore-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type bookHandler struct {
	service service.Service[model.Book]
}

func NewBookHandler(service service.Service[model.Book]) *bookHandler {
	return &bookHandler{service}
}

func (h *bookHandler) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		return ctx.Status(400).JSON(nil)
	}

	user, err := h.service.Find(id)
	if err != nil {
		return ctx.Status(404).JSON(nil)
	}

	return ctx.JSON(&user)
}

func (h *bookHandler) GetAll(ctx *fiber.Ctx) error {
	users, err := h.service.FindAll()
	if err != nil {
		return ctx.Status(404).JSON(nil)
	}

	return ctx.JSON(&users)
}
