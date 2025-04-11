package main

import (
	"log"
	"os"

	"come-back/controller"
	"come-back/middleware"
	"come-back/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading environment variables!")
	}

	if err := repository.InitMySQL(os.Getenv("MYSQL_DSN")); err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	api := router.Group("/api")
	{
		public := api.Group("")
		{

			public.GET("/posts", controller.GetPostsPaginated)
			public.GET("/post/:id", controller.GetPost)
			public.GET("/post/:id/comments", controller.GetPostComments)

		}

		auth := api.Group("").Use(middleware.UserAuth())
		{
			auth.POST("/post", controller.CreatePost)
			auth.POST("/post/:id/comment", controller.CreateComment)
			auth.PUT("/post/:id", controller.UpdatePost)
			auth.DELETE("/post/:id", controller.DeletePost)
		}

		admin := api.Group("/admin").Use(middleware.AdminAuth())
		{
			admin.GET("/users", controller.GetAllUsers)
			admin.PUT("/users/:id/ban", controller.BanUser)
			admin.PUT("/users/:id/promote", controller.PromoteToAdmin)
			admin.DELETE("/post/:id", controller.DeletePostAdmin)
			admin.DELETE("/comments/:id", controller.DeleteCommentAdmin)
			admin.GET("/dashboard", controller.AdminDashboard)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
