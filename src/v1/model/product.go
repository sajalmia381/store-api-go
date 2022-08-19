package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID  `json:"id" bson:"_id"`
	CreatedBy   string              `json:"createdBy" bson:"createdBy"`
	Category    *primitive.ObjectID `json:"category" bson:"category"`
	ImageSource *primitive.ObjectID `json:"imageSource" bson:"imageSource"`
	Title       string              `json:"title" bson:"title"`
	Slug        string              `json:"slug" bson:"slug"`
	Price       int                 `json:"price" bson:"price"`
	Image       string              `json:"image" bson:"image"`
	Description string              `json:"description" bson:"description"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
	Active      bool                `json:"active" bson:"active"`
}
