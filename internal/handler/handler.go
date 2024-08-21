package handler

import "github.com/gofiber/fiber/v2"

type Handler interface {
	Get(ctx *fiber.Ctx)
	Post(ctx *fiber.Ctx)
	GetAll(ctx *fiber.Ctx)
}
