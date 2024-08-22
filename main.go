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
		userRepo     = repository.NewUserRepo[model.User](database.DBConn)
		bookRepo     = repository.NewBookRepo[model.Book](database.DBConn)
		orderRepo    = repository.NewOrderRepo[model.Order](database.DBConn)
		userService  = service.NewService(userRepo)
		bookService  = service.NewService(bookRepo)
		orderService = service.NewService(orderRepo)

		userHandler  = handler.NewUserHandler(userService, orderService)
		authHandler  = handler.NewAuthHandler(userService)
		bookHandler  = handler.NewBookHandler(bookService)
		orderHandler = handler.NewOrderHandler(orderService, userService)

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

	// user routes
	userRoutes := apiv1.Group("/users")
	userRoutes.Post("/", userHandler.Post)
	userRoutes.Get("/me/orders", middleware.JWTAuthentication, userHandler.GetOrders)

	// book routes
	bookRoutes := apiv1.Group("/books", middleware.JWTAuthentication)
	bookRoutes.Get("/", bookHandler.GetAll)

	// order routes
	orderRoutes := apiv1.Group("/orders", middleware.JWTAuthentication)
	orderRoutes.Post("/", orderHandler.Post)

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Panic(err)
		}
	}()

	// Create channel to signify a signal being sent
	c := make(chan os.Signal, 1)
	// When an interrupt or termination signal is sent, notify the channel
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

}
