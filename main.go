package main

import (
	"os"
	"phuhung273/progress-tracking/db"
	"phuhung273/progress-tracking/router"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	db.Init()
	
	router.Init()

	router.Route()

	port := os.Getenv("PORT")
	router.Router.Listen(":" + port)

}