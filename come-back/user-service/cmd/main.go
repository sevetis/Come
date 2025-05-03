package main

import (
	"log"
	"os"
	"user-service/internal/controller"
	"user-service/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	if err := repository.InitDB(os.Getenv("MYSQL_DSN")); err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.Static("/uploads", "./uploads")

	api := router.Group("/api/user")
	{
		public := api.Group("")
		{
			public.POST("/login", controller.Login)
			public.POST("/register", controller.Register)

			public.GET("/:id", controller.GetUser)
			public.GET("/batch", controller.GetUsersBatch)
		}

		auth := api.Group("").Use(controller.UserAuth())
		{
			auth.GET("/profile", controller.GetProfile)
			auth.PUT("/profile", controller.UpdateProfile)
			auth.POST("/avatar", controller.UploadAvatar)
		}

		admin := api.Group("/admin").Use(controller.AdminAuth())
		{
			admin.GET("/users", controller.GetAllUsers)
			admin.PUT("/users/:id/ban", controller.BanUser)
			admin.PUT("/users/:id/promote", controller.PromoteToAdmin)
		}
	}

	port := os.Getenv("USER_PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("User service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start user service: %v", err)
	}
}
