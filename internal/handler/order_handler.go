package handler

import (
	"bookstore-api/internal/model"
	"bookstore-api/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type orderHandler struct {
	service     service.Service[model.Order]
	userService service.Service[model.User]
}

func NewOrderHandler(service service.Service[model.Order], userService service.Service[model.User]) *orderHandler {
	return &orderHandler{service, userService}
}

func (h *orderHandler) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	if err != nil {
		return ctx.Status(400).JSON(nil)
	}

	order, err := h.service.Find(id)
	if err != nil {
		return ctx.Status(404).JSON(nil)
	}

	return ctx.JSON(&order)
}

func (h *orderHandler) Post(ctx *fiber.Ctx) error {

	orderReq := orderRequest{}

	userEmail, ok := ctx.Locals("userEmail").(string)
	if !ok || userEmail == "" {
		return ctx.SendStatus(fiber.StatusForbidden)
	}

	if err := ctx.BodyParser(&orderReq); err != nil {
		return ctx.Status(500).JSON(nil)
	}

	validate := validator.New()
	if err := validate.Struct(orderReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	user, err := h.userService.FindByCondition("email = ?", userEmail)
	if err != nil {

		return ctx.SendStatus(fiber.StatusForbidden)
	}

	order := NewOrderFromRequest(orderReq)

	order.UserID = user.ID

	u, err := h.service.Create(*order)
	if err != nil {
		return ctx.Status(500).JSON(nil)
	}
	log.Info("order created with id", u.ID)

	return ctx.JSON(&u)

}

func (h *orderHandler) GetAll(ctx *fiber.Ctx) error {
	orders, err := h.service.FindAll()
	if err != nil {
		return ctx.Status(404).JSON(nil)
	}

	return ctx.JSON(&orders)
}
