package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-playground/db"
	"go-playground/logger"
	"log"
	"os"
)

func main() {
	Initialize()

	e := echo.New()

	// Logger()는 서버로 들어오는 모든 요청을 로그로 남겨줍니다.
	e.Use(middleware.Logger())

	// Recover()는 서버로 요청이 들어오고, 해당 요청으로 인해 golang 서버에 panic 이 발생했을 경우에도 golang 서버가 죽지않고 동작할 수 있도록 회복시켜주는 역할을 합니다.
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(os.Getenv("PORT")))
}

func Initialize() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database connection
	db.Connect()

	// Initialize the logger
	logger.InitLogger()
}
