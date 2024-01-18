package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	Id         primitive.ObjectID `bson:"_id"`
	Name       *string            `json:"name" validate:"required"`
	Topic      *string            `json:"topic" validate:"required"`
	Genre_id   int                `json:"genre_id"`
	Movie_URL  *string            `json:"movie_url" validate:"required"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	Movie_id   int                `json:"movie_id"`
}
