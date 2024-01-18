package routes

import (
	"arpithku/movie_app/controllers"
	"arpithku/movie_app/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router gin.Engine) {
	router.Use(middleware.AuthenticateUser())
	router.GET("/users/:user_id", controllers.GetUser())
}
