package db

import (
	"fmt"
	"os"

	"phuhung273/progress-tracking/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s", dbHost, dbName, dbUser, dbPwd)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
    	panic("failed to connect database")
  	}

	// Migrate the schema
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Exercise{})
	DB.AutoMigrate(&models.Setting{})
	DB.AutoMigrate(&models.Result{})
}