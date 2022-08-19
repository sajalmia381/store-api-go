package model

import (
	"time"

	"github.com/sajalmia381/store-api/src/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password" bson:"password"`
	Number      *int               `json:"number" bson:"number"`
	Status      bool               `json:"status" bson:"status"`
	Role        enums.Role         `json:"role" bson:"role"`
	LastLoginAt *time.Time         `json:"lastLoginAt" bson:"lastLoginAt"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updateAt" bson:"updateAt"`
}
