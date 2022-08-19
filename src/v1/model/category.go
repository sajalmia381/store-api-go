package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID          primitive.ObjectID  `json:"id" bson:"_id"`
	Parent      *primitive.ObjectID `json:"parent" bson:"parent"`
	Products    []primitive.ObjectID
	Name        string `json:"name" bson:"name"`
	Slug        string `json:"slug" bson:"slug"`
	Description string `json:"description" bson:"description"`
}
