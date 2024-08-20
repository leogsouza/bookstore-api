package main

import (
	"bookstore-api/internal/database"
	"bookstore-api/internal/model"
	"log"
)

func main() {
	books := []model.Book{
		{
			Title:  "The Invisible Man",
			Author: "H. G. Wells",
			Price:  5.75,
		},
		{
			Title:  "The War of the Worlds",
			Author: "H. G. Wells",
			Price:  8.50,
		},
		{
			Title:  "1984",
			Author: "George Orwell",
			Price:  8.50,
		},
		{
			Title:  "War and Peace",
			Author: "Leo Tolstoy",
			Price:  12.25,
		},
		{
			Title:  "Moby Dick",
			Author: "Hearman Melville",
			Price:  10.99,
		},
		{
			Title:  "Dracula",
			Author: "Bram Stocker",
			Price:  7.39,
		},
		{
			Title:  "Hamlet",
			Author: "William Shakespeare",
			Price:  20.25,
		},
		{
			Title:  "The Last of the Mohicans",
			Author: "James Fenimore Cooper",
			Price:  15.50,
		},
		{
			Title:  "The Arabian Nights",
			Author: "Andrew Lang",
			Price:  17.80,
		},
		{
			Title:  "Journey to the Center of the Earth",
			Author: "Jules Verne",
			Price:  10.88,
		},
	}

	database.ConnectDB()

	result := database.DBConn.Create(books)
	if result.Error != nil {
		log.Fatalf("an error occured when seeding the database: %v", result.Error)
	}

	log.Printf("Books created!")
}
