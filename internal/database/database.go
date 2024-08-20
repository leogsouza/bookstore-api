package database

import (
	"bookstore-api/internal/model"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbHost     = os.Getenv("DB_HOST")
	dbPort     = os.Getenv("DB_PORT")
	dbDatabase = os.Getenv("DB_DATABASE")
	dbUsername = os.Getenv("DB_USERNAME")
	dbPassword = os.Getenv("DB_PASSWORD")
)

var DBConn *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/New_York",
		dbHost, dbUsername, dbPassword, dbDatabase, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	log.Println("connected")
	log.Println("running migrations")
	db.AutoMigrate(&model.Book{}, &model.User{}, &model.Order{}, &model.OrderItem{})

	DBConn = db
}
