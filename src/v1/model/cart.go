package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartProductSpec struct {
	ProductId primitive.ObjectID `json:"productId" bson:"productId"`
	Quantity  uint16             `json:"quantity" bson:"quantity"`
}

type Cart struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserId    primitive.ObjectID `json:"userId" bson:"userId"`
	Products  []CartProductSpec  `json:"products" bson:"products"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
