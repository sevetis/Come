package main

import (
	"chat-service/internal/client"
	"chat-service/internal/controller"
	"chat-service/internal/repository"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	repository.InitDB(os.Getenv("MYSQL_DSN"))
	controller.InitRedis()

	chatCtrl := controller.NewChatController(
		*client.NewUserClient(fmt.Sprintf("http://user-service:%s", os.Getenv("USER_PORT"))),
	)

	router := gin.Default()
	router.GET("/api/history", chatCtrl.GetChatHistory)
	router.GET("/api/chatroom", controller.WsOptionalAuth(), chatCtrl.HandleChat)
	router.GET("/api/online", chatCtrl.GetOnlineCount)

	port := os.Getenv("CHAT_PORT")
	if port == "" {
		port = "8083"
	}
	log.Printf("Chat service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start chat service: %v", err)
	}
}
