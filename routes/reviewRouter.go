package routes

import (
	controllers "arpithku/movie_app/controllers"
	"arpithku/movie_app/middleware"

	"github.com/gin-gonic/gin"
)

func ReviewRoutes(router gin.Engine) {
	router.Use(middleware.Authenticate())
	router.POST("reviews/addreview", controllers.AddAReview())
}
