package main

import (
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
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})



	r.Run()
}
