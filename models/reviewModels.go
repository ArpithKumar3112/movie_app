package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reviews struct {
	Id          primitive.ObjectID `bson:"_id"`
	Movie_id    int                `json:"movie_id"`
	Review_id   int                `json:"review_id"`
	Review      *string            `json:"review" validate:"required"`
	Created_at  time.Time          `json:"created_at"`
	Updated_at  time.Time          `json:"updated_at"`
	Reviewer_id string             `json:"reviewer_id"`
}
