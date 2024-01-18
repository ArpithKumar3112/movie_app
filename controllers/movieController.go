package controllers

import (
	"arpithku/movie_app/database"
	helper "arpithku/movie_app/helper"
	"arpithku/movie_app/models"
	"context"
	"log"
	"net/http"
	"time"

	//"strconv"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection(database.Client, "movie")

// To create one movie
func CreateMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.VerifyUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var movie models.Movie
		defer cancel()

		if err := c.BindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		//Check to see if name exists
		regexMatch := bson.M{"$regex": primitive.Regex{Pattern: *movie.Name, Options: "i"}}
		count, err := movieCollection.CountDocuments(ctx, bson.M{"name": regexMatch})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while checking for the movie name"})
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this movie name already exists", "count": count})
			return
		}

		if validationError := validate.Struct(&movie); validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": validationError.Error()}})
			return
		}

		newMovie := models.Movie{
			Id:       primitive.NewObjectID(),
			Name:     movie.Name,
			Topic:    movie.Topic,
			Genre_id: movie.Genre_id,

			Movie_URL:  movie.Movie_URL,
			Created_at: movie.Created_at,
			Updated_at: movie.Updated_at,
			Movie_id:   movie.Movie_id,
		}
		result, err := movieCollection.InsertOne(ctx, newMovie)
		//err = movieCollection.FindOne(ctx, bson.M{"movie_id": movie.Movie_id}).Err()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"Status":  http.StatusCreated,
			"Message": "success",
			"Data":    map[string]interface{}{"data": result}})
	}
}
