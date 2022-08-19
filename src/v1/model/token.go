package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	ID        primitive.ObjectID  `json:"id" bson:"_id"`
	UserId    *primitive.ObjectID `json:"userId" bson:"userId"`
	Token     string              `json:"token" bson:"token"`
	Type      string              `json:"type" bson:"type"`
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt"`
}
