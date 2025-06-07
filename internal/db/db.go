package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	// LOAD ENV AND CATCH ERROR
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error in loading .env file")
	}
}

func InitDB() (*gorm.DB, *sql.DB) {

	// GET ENV VARIABLES FROM ENV
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// CONNECTION STRING - DATA SOURCE NAME TO ESTABLISH CONNECTION
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	// CONNECT TO DB
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database : ", err)
		return nil, nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("Unable to get sql.DB from gorm.DB:", err)
	}

	fmt.Println("Successfully connected to PostgreSQL ...")

	return db, sqlDB

}
