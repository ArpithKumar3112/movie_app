package routes

import (
	controllers "arpithku/movie_app/controllers"
	"arpithku/movie_app/middleware"

	"github.com/gin-gonic/gin"
)

func MovieRoutes(router gin.Engine) {
	router.Use(middleware.AuthenticateUser())
	router.POST("/movies/createmovie", controllers.CreateMovie())
}
