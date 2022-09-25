package main

import (
	"os"
	"phuhung273/progress-tracking/db"
	"phuhung273/progress-tracking/router"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db.Init()
	
	router.Init()

	router.Route()

	port := os.Getenv("PORT")
	router.Router.Listen(":" + port)

}