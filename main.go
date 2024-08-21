package main

import (
	"bookstore-api/internal/database"
	"bookstore-api/internal/handler"
	"bookstore-api/internal/model"
	"bookstore-api/internal/repository"
	"bookstore-api/internal/service"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	database.ConnectDB()
	app := fiber.New(fiber.Config{
		BodyLimit: 200 * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Error("Error happened ", err)
			c.Status(fiber.StatusInternalServerError)
			return c.SendString(err.Error())
		},
	})

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip}  ${status} - ${latency} ${method} ${path}\n",
	}))

	app.Use(recover.New())
	app.Use(cors.New())

	api := app.Group("/api")

	v1 := api.Group("/v1")

	userRepo := repository.NewUserRepo[model.User](database.DBConn)
	userService := service.NewService(userRepo)
	userRoutes := v1.Group("/users")

	userHandler := handler.NewUserHandler(userService)
	userRoutes.Get("/:id", userHandler.Get)
	userRoutes.Get("/", userHandler.GetAll)
	userRoutes.Post("/", userHandler.Post)

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

}
