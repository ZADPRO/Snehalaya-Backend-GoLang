package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error in loading .env file")
	}
}

func InitDB() (*gorm.DB, *sql.DB) {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	var dbConn *gorm.DB
	var sqlDB *sql.DB

	maxRetries := 10

	for i := 1; i <= maxRetries; i++ {
		fmt.Printf("⏳ Connecting to PostgreSQL... attempt %d/%d\n", i, maxRetries)

		dbTemp, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			fmt.Println("✅ Successfully connected to PostgreSQL!")

			sqlTemp, err := dbTemp.DB()
			if err == nil {
				dbConn = dbTemp
				sqlDB = sqlTemp
				return dbConn, sqlDB
			}
		}

		fmt.Println("❌ DB connection failed. Retrying in 3 seconds...")
		time.Sleep(3 * time.Second)
	}

	fmt.Println("🚨 Could not connect to DB after retries. Server will still start.")
	return nil, nil
}
