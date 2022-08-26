package dtos

import (
	"crypto/rsa"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JwtPayload struct {
	ID    primitive.ObjectID `json:"id" bson:"_id" header:"id"`
	Name  string             `json:"name" bson:"name" header:"name"`
	Email string             `json:"email" bson:"email" header:"email"`
	Role  string             `json:"role" bson:"role" header:"role"`
}

type JwtUserData struct {
	ID    string `json:"id" bson:"_id" header:"id"`
	Name  string `json:"name" bson:"name" header:"name"`
	Email string `json:"email" bson:"email" header:"email"`
	Role  string `json:"role" bson:"role" header:"role"`
}

type JwtResponseDto struct {
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}

type RsaKeys struct {
	TokenSecretKey        *rsa.PublicKey
	RefreshTokenSecretKey *rsa.PrivateKey
}
