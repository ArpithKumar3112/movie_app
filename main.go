package main

import (
	"arpithku/movie_app/database"
	"fmt"
	"os"

	//"arpithku/movie_app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Welcome to the movie_app")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router := gin.Default()

	//run Database
	database.StartDB()

	//Log events
	router.Use(gin.Logger())

	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": "Welcome to movie_app api!",
		})
	})

	router.Run(":" + port)
}
