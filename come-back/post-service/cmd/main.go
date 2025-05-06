package main

import (
	"log"
	"os"
	"post-service/internal/controller"
	"post-service/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading environment variables!")
	}

	if err := repository.InitDB(os.Getenv("MYSQL_DSN")); err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	api := router.Group("/api/post")
	{
		public := api.Group("")
		{
			public.GET("/batch", controller.GetPostsPaginated)
			public.GET("/:id", controller.GetPost)
			public.GET("/:id/comments", controller.GetPostComments)

		}

		auth := api.Group("").Use(controller.UserAuth())
		{
			auth.POST("/create", controller.CreatePost)
			auth.POST("/:id/comment", controller.CreateComment)
			auth.PUT("/:id", controller.UpdatePost)
			auth.DELETE("/:id", controller.DeletePost)
		}

		admin := api.Group("/admin").Use(controller.AdminAuth())
		{
			admin.DELETE("/:id", controller.DeletePostAdmin)
			admin.DELETE("/comments/:id", controller.DeleteCommentAdmin)
		}
	}

	port := os.Getenv("POST_PORT")
	if port == "" {
		port = "8082"
	}
	router.Run(":" + port)
}
