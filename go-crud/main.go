package main

import (
	"go-crud/controllers"
	"go-crud/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	// Load environment variables from .env file
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	r.POST("/posts", controllers.PostsCreate)

	r.PUT("/posts/:id", controllers.PostsUpdate)

	r.GET("/posts", controllers.PostsIndex)

	r.GET("/posts/:id", controllers.PostsShow)

	r.DELETE("/posts/:id", controllers.PostsDelete)

	r.Run()
}
