package handler

import (
	"bookstore-api/internal/model"
	"bookstore-api/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type userHandler struct {
	service service.Service[model.User]
}

func NewUserHandler(service service.Service[model.User]) *userHandler {
	return &userHandler{service}
}

func (h *userHandler) Get(ctx *fiber.Ctx) error {
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

func (h *userHandler) Post(ctx *fiber.Ctx) error {

	userReq := userRequest{}

	u := model.User{}
	if err := ctx.BodyParser(&userReq); err != nil {
		return ctx.Status(500).JSON(nil)
	}

	validate := validator.New()
	if err := validate.Struct(userReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	u.Name = userReq.Name
	u.Email = userReq.Email
	u.Password = userReq.Password

	userID, err := h.service.Create(u)
	if err != nil {
		return ctx.Status(500).JSON(nil)
	}
	log.Info("user created with id", userID)

	return ctx.JSON(&u)

}

func (h *userHandler) GetAll(ctx *fiber.Ctx) error {
	users, err := h.service.FindAll()
	if err != nil {
		return ctx.Status(404).JSON(nil)
	}

	return ctx.JSON(&users)
}
