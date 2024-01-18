package controllers

import (
	"arpithku/movie_app/database"
	"arpithku/movie_app/models"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	helper "arpithku/movie_app/helper"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var genreCollection *mongo.Collection = database.OpenCollection(database.Client, "genre")
var validate = validator.New()

func CreateGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.VerifyUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var genre models.Genre
		defer cancel()

		if err := c.BindJSON(&genre); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		regexMatch := bson.M{"$regex": primitive.Regex{Pattern: *genre.Name, Options: "i"}}
		count, err := genreCollection.CountDocuments(ctx, bson.M{"Name": regexMatch})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "error occured while checking for the genre name"})
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this genre name already exists", "count": count})
			return
		}
		if validationError := validate.Struct(&genre); validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": validationError.Error()}})
			return
		}
		newGenre := models.Genre{
			Id:         primitive.NewObjectID(),
			Name:       genre.Name,
			Created_at: time.Now(),
			Updated_at: time.Now(),
			Genre_id:   genre.Genre_id, // Insert genre_id
		}
		result, err := genreCollection.InsertOne(ctx, newGenre)
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
			"Message": "genre created successfully",
			"Data":    map[string]interface{}{"data": result}})
	}
}

func GetGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		genre_id := c.Param("genre_id")
		var genre models.Genre
		defer cancel()

		genreId, erro := strconv.Atoi(genre_id)
		if erro != nil {
			//Handle error
		}
		err := genreCollection.FindOne(ctx, bson.M{"genre_id": genreId}).Decode(&genre)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "error",
				"Data":    map[string]interface{}{"error": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Status":  http.StatusOK,
			"Message": "success",
			"Data":    map[string]interface{}{"data": genre}})
	}
}

func GetGenres() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "genre_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}}}}
		result, err := genreCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching genre details"})
		}
		var allgenres []bson.M
		if err = result.All(ctx, &allgenres); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allgenres[0])
	}
}

func EditGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.VerifyUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		genreId := c.Param("genreId")
		var genre models.Genre
		defer cancel()
		i, erro := strconv.Atoi(genreId)
		if erro != nil {
			//Handle the error
		}
		if err := c.BindJSON(&genre); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}
		if validationErr := validate.Struct(&genre); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status":  http.StatusBadRequest,
				"message": "error",
				"Data":    map[string]interface{}{"data": validationErr.Error()}})
			return
		}
		update := bson.M{"name": genre.Name}
		filterbyId := bson.M{"genre_id": i}

		result, err := genreCollection.UpdateOne(ctx, filterbyId, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}
		var updatedGenre models.Genre
		if result.MatchedCount == 1 {
			err := genreCollection.FindOne(ctx, filterbyId).Decode(&updatedGenre)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"Status":  http.StatusInternalServerError,
					"Message": "error",
					"Data":    map[string]interface{}{"data": err.Error()}})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"Status":  http.StatusOK,
			"Message": "success",
			"Data":    updatedGenre})
	}
}

func DeleteGenre() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		genreId := c.Param("genre_id")
		defer cancel()
		i, erro := strconv.Atoi(genreId)
		if erro != nil {
			//Handle the error
		}

		result, err := genreCollection.DeleteOne(ctx, bson.M{"genre_id": i})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				gin.H{
					" Status":  http.StatusNotFound,
					" Message": "error",
					" Data":    map[string]interface{}{"data": "Genre with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			gin.H{
				"Status":  http.StatusOK,
				"Message": "success",
				"Data":    map[string]interface{}{"data": "Genre successfully deleted!"}},
		)
	}
}

func SearchByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var genre models.Genre
		defer cancel()
		// if err:=c.BindJSON(&genre);err!=nil{
		// 	c.JSON(http.StatusBadRequest,gin.H{
		// 		"Status":http.StatusBadRequest,
		// 		"Message":"error",
		// 		"Data":map[string]interface{}{
		// 			"data":err.Error()}})
		// 		return
		// 	}
		//name:=genre.Name
		err := genreCollection.FindOne(ctx, bson.M{"name": "Animation"}).Decode(&genre)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status":  http.StatusInternalServerError,
				"Message": "error",
				"Data":    map[string]interface{}{"error": err.Error()},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Status":  http.StatusOK,
			"Message": "success",
			"Data":    map[string]interface{}{"data": genre}})
	}
}
