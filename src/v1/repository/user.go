package repository

import (
	"errors"
	"log"
	"time"

	"github.com/sajalmia381/store-api/src/enums"
	"github.com/sajalmia381/store-api/src/v1/db"
	"github.com/sajalmia381/store-api/src/v1/dtos"
	"github.com/sajalmia381/store-api/src/v1/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	Store(user model.User) (model.User, error)
	FindAll(queryParams dtos.UserQuery) ([]model.User, error)
	FindById(id primitive.ObjectID) (model.User, error)
	FindByEmail(email string) (model.User, error)
	UpdateById(id primitive.ObjectID, payload primitive.M) (model.User, error)
	UpdateLoginTime(id primitive.ObjectID) (model.User, error)
	DeleteById(id primitive.ObjectID) (*mongo.DeleteResult, error)
}

type userRepository struct {
	dm *db.DmManager
}

func (r userRepository) Store(user model.User) (model.User, error) {
	coll := r.dm.DB.Collection(string(enums.USER_COLLECTION_NAME))
	_, err := coll.InsertOne(r.dm.Ctx, user)
	if err != nil {
		log.Println("[ERROR] Insert document:", err.Error())
		// log.Printf("type %T", err)
		return user, err
	}
	return user, nil
}

func (r userRepository) FindAll(filterData dtos.UserQuery) ([]model.User, error) {
	var objects []model.User
	query := bson.D{}
	if filterData.Role != "" {
		query = append(query, bson.E{
			Key: "role", Value: filterData.Role,
		})
	}
	if filterData.Status != nil {
		query = append(query, bson.E{
			Key: "status", Value: *filterData.Status,
		})
	}
	coll := r.dm.DB.Collection(string(enums.USER_COLLECTION_NAME))
	cursor, err := coll.Find(r.dm.Ctx, query)
	if err != nil {
		return objects, err
	}
	if err := cursor.All(r.dm.Ctx, &objects); err != nil {
		log.Println("[ERROR]:", err.Error())
		panic(err)
	}
	return objects, nil
}

func (r userRepository) FindById(id primitive.ObjectID) (model.User, error) {
	var user model.User
	query := bson.D{
		{Key: "_id", Value: id},
	}
	coll := r.dm.DB.Collection(string(enums.USER_COLLECTION_NAME))
	result := coll.FindOne(r.dm.Ctx, query)
	if err := result.Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return user, errors.New("user is not found")
		}
		log.Println("[ERROR]:", err.Error())
		panic(err)
	}
	return user, nil
}

func (r userRepository) FindByEmail(email string) (model.User, error) {
	var user model.User
	filter := bson.D{
		{Key: "email", Value: email},
	}
	coll := r.dm.DB.Collection(string(enums.USER_COLLECTION_NAME))
	result := coll.FindOne(r.dm.Ctx, filter)

	if err := result.Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return user, errors.New("user is not found")
		}
		log.Println("[ERROR]:", err.Error())
		panic(err)
	}

	return user, nil
}

func (r userRepository) UpdateById(id primitive.ObjectID, payload primitive.M) (model.User, error) {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	payload["updatedAt"] = time.Now().UTC()
	filter := bson.D{
		{Key: "_id", Value: id},
	}
	update := bson.D{
		{Key: "$set", Value: payload},
	}
	coll := r.dm.DB.Collection(string(enums.USER_COLLECTION_NAME))
	result := coll.FindOneAndUpdate(r.dm.Ctx, filter, update, opts)
	var user model.User
	if err := result.Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return user, err
		}
		log.Println("[ERROR] Update document:", err.Error())
		panic(err)
	}

	return user, nil
}

func (r userRepository) UpdateLoginTime(id primitive.ObjectID) (model.User, error) {
	var user model.User
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	filter := bson.D{
		{Key: "_id", Value: id},
	}
	update := bson.D{
		{Key: "$set", Value: bson.M{
			"lastLoginAt": time.Now().UTC(),
		}},
	}
	coll := r.dm.DB.Collection(string(enums.USER_COLLECTION_NAME))
	result := coll.FindOneAndUpdate(r.dm.Ctx, filter, update, opts)
	if err := result.Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return user, err
		}
		log.Println("[ERROR] Update Login time document:", err.Error())
		panic(err)
	}
	return user, nil
}

func (r userRepository) DeleteById(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	query := bson.D{
		{Key: "_id", Value: id},
	}

	coll := r.dm.DB.Collection(string(enums.USER_COLLECTION_NAME))
	result, err := coll.DeleteOne(r.dm.Ctx, query)
	if err != nil {
		return nil, err
	}
	if result.DeletedCount != 1 {
		log.Println("[ERROR] Delete document count", result.DeletedCount)
		return result, errors.New("user is not exists")
	}
	return result, nil
}

func NewUserRepository() UserRepository {
	return &userRepository{
		dm: db.GetDmManager(),
	}
}
