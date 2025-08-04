package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"time"
)

var database *gorm.DB

func Connect() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌  Failed to connect to the database: %v", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("❌  Failed to get database instance: %v", err)
	}

	maxConn, err := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	if err != nil {
		maxConn = 20
		log.Printf("⚠️ Failed to parse DB_MAX_CONNECTIONS environment variable, using default value [%d].", maxConn)
	}
	idleConn, err := strconv.Atoi(os.Getenv("DB_IDLE_CONNECTIONS"))
	if err != nil {
		idleConn = 10
		log.Printf("⚠️ Failed to parse DB_IDLE_CONNECTIONS environment variable, using default value [%d].", idleConn)
	}

	sqlDB.SetMaxOpenConns(maxConn)
	sqlDB.SetMaxIdleConns(idleConn)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	log.Println("✅  Connected to database successfully.")
	log.Println("⚙️ HOST =", os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT"))
	log.Println("⚙️ MAX_CONNECTION =", maxConn)
	log.Println("⚙️ IDLE_CONNECTION =", idleConn)
	log.Println()

	return database
}
