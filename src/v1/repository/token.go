package repository

import (
	"errors"
	"time"

	"github.com/sajalmia381/store-api/src/enums"
	"github.com/sajalmia381/store-api/src/v1/db"
	"github.com/sajalmia381/store-api/src/v1/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	TokenCollectionName = "tokens"
)

type TokenRepository interface {
	Store(payload model.Token) (model.Token, error)
	FindByToken(token string) (model.Token, error)
	DeleteByToken(token string) (*mongo.DeleteResult, error)
}

type tokenRepository struct {
	dm *db.DmManager
}

func (r tokenRepository) Store(payload model.Token) (model.Token, error) {
	payload.CreatedAt = time.Now().UTC()
	if payload.Type == "" {
		payload.Type = string(enums.REFRESH_TOKEN)
	}
	coll := r.dm.DB.Collection(TokenCollectionName)
	_, err := coll.InsertOne(r.dm.Ctx, payload)
	if err != nil {
		return payload, err
	}
	return payload, nil
}

func (r tokenRepository) FindByToken(token string) (model.Token, error) {
	var object model.Token
	filter := bson.D{
		{Key: "token", Value: token},
	}
	coll := r.dm.DB.Collection(TokenCollectionName)
	result := coll.FindOne(r.dm.Ctx, filter)
	if err := result.Decode(&object); err != nil {
		if err == mongo.ErrNoDocuments {
			return object, errors.New("token object is not found")
		}
		panic(err)
	}
	return object, nil
}

func (r tokenRepository) DeleteByToken(token string) (*mongo.DeleteResult, error) {
	filter := bson.D{
		{Key: "token", Value: token},
	}
	coll := r.dm.DB.Collection(TokenCollectionName)
	result, err := coll.DeleteOne(r.dm.Ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, errors.New("token object is not found")
		}
		panic(err)
	}
	return result, nil
}

func NewTokenRepository() TokenRepository {
	return &tokenRepository{
		dm: db.GetDmManager(),
	}
}
