package handler

import (
	"bookstore-api/internal/model"
	"bookstore-api/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type userHandler struct {
	service      service.Service[model.User]
	orderService service.Service[model.Order]
}

func NewUserHandler(service service.Service[model.User], orderService service.Service[model.Order]) *userHandler {
	return &userHandler{service, orderService}
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

	if err := ctx.BodyParser(&userReq); err != nil {
		return ctx.Status(500).JSON(nil)
	}

	validate := validator.New()
	if err := validate.Struct(userReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	user, err := NewUserFromRequest(userReq)
	if err != nil {
		return err
	}
	u, err := h.service.Create(*user)
	if err != nil {
		return ctx.Status(500).JSON(nil)
	}
	log.Info("user created with id", u.ID)

	return ctx.JSON(&u)

}

func (h *userHandler) GetAll(ctx *fiber.Ctx) error {
	users, err := h.service.FindAll()
	if err != nil {
		return ctx.Status(404).JSON(nil)
	}

	return ctx.JSON(&users)
}

func (h *userHandler) GetOrders(ctx *fiber.Ctx) error {
	userEmail, ok := ctx.Locals("userEmail").(string)
	if !ok || userEmail == "" {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	user, err := h.service.FindByCondition("email = ?", userEmail)
	if err != nil {

		return ctx.SendStatus(fiber.StatusForbidden)
	}
	orders, err := h.orderService.FindAllByCondition("user_id = ?", user.ID)
	if err != nil {
		return ctx.Status(500).JSON(nil)
	}

	return ctx.JSON(&orders)
}
