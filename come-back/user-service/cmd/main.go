package main

import (
	"come-back/user-service/internal/controller"
	"come-back/user-service/internal/repository"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	repository.InitDB(os.Getenv("MYSQL_DSN"))

	router := gin.Default()
	router.POST("/api/login", controller.Login)
	router.POST("/api/register", controller.Register)
	router.GET("/api/users/:id", controller.GetUser)
	router.GET("/api/users/batch", controller.GetUsersBatch)

	port := os.Getenv("USER_PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("User service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start user service: %v", err)
	}
}
