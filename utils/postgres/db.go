package util

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"testjun/models"
)

func GetDbInstance() *gorm.DB {
	godotenv.Load(".env")

	dbHost := os.Getenv("HOST")
	dbUser := os.Getenv("DBUSER")
	dbPassword := os.Getenv("DBPASSWORD")
	dbName := os.Getenv("DBNAME")
	dbPort := os.Getenv("PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow", dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}
	return db.Debug()
}

func AutoMigrate(connection *gorm.DB) {
	connection.Debug().AutoMigrate(
		&models.User{})
}
