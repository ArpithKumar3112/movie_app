package routes

import (
	"arpithku/movie_app/controllers"
	"arpithku/movie_app/middleware"

	"github.com/gin-gonic/gin"
)

func GenreRoutes(incomingRoutes gin.Engine) {
	incomingRoutes.Use(middleware.AuthenticateUser())
	incomingRoutes.POST("/genres/creategenre", controllers.CreateGenre())
	incomingRoutes.GET("/genres/:genre_id", controllers.GetGenre())
	incomingRoutes.GET("/genres", controllers.GetGenres())
	incomingRoutes.PUT("/genres/:genre_id", controllers.EditGenre())
	incomingRoutes.DELETE("/genres/:genre_id", controllers.DeleteGenre())
}
