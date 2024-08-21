package main

import (
	"bookstore-api/internal/database"
	"bookstore-api/internal/handler"
	"bookstore-api/internal/middleware"
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

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		log.Error("Error happened ", err)
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	database.ConnectDB()

	var (
		userRepo    = repository.NewUserRepo[model.User](database.DBConn)
		userService = service.NewService(userRepo)

		userHandler = handler.NewUserHandler(userService)
		authHandler = handler.NewAuthHandler(userService)

		app = fiber.New(config)

		apiv1 = app.Group("/api/v1")
	)

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip}  ${status} - ${latency} ${method} ${path}\n",
	}))

	app.Use(recover.New())
	app.Use(cors.New())

	// auth routes
	auth := apiv1.Group("/auth")
	auth.Post("/login", authHandler.Authenticate)

	userRoutes := apiv1.Group("/users", middleware.JWTAuthentication)

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
