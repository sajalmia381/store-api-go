package repository

import (
	"errors"
	"log"
	"time"

	"github.com/sajalmia381/store-api/src/v1/db"
	"github.com/sajalmia381/store-api/src/v1/dtos"
	"github.com/sajalmia381/store-api/src/v1/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Collection Name
const (
	UserCollectionName = "users"
)

// Error
var (
	invalidIdErr = errors.New("invalid user id")
)

type UserRepository interface {
	Store(user model.User) (model.User, error)
	FindAll(queryParams dtos.UserQuery) ([]model.User, error)
	FindById(id string) (model.User, error)
	FindByEmail(email string) (model.User, error)
	UpdateById(id string, payload dtos.UserUpdateDto) (model.User, error)
	UpdateLoginTime(id string) (model.User, error)
	DeleteById(id string) (*mongo.DeleteResult, error)
}

type userRepository struct {
	dm *db.DmManager
}

func (r userRepository) Store(user model.User) (model.User, error) {
	user.ID = primitive.NewObjectID()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[ERROR] convert string to hash password:", err.Error())
		return user, err
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()

	coll := r.dm.DB.Collection(UserCollectionName)
	_, err = coll.InsertOne(r.dm.Ctx, user)
	if err != nil {
		log.Println("[ERROR] Insert document:", err.Error())
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
	coll := r.dm.DB.Collection(UserCollectionName)
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

func (r userRepository) FindById(id string) (model.User, error) {
	user := model.User{}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, invalidIdErr
	}
	query := bson.D{
		{Key: "_id", Value: _id},
	}
	coll := r.dm.DB.Collection(UserCollectionName)
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
	coll := r.dm.DB.Collection(UserCollectionName)
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

func (r userRepository) UpdateById(id string, payload dtos.UserUpdateDto) (model.User, error) {
	user, err := r.FindById(id)
	if err != nil {
		return user, err
	}
	payload.UpdateAt = time.Now().UTC()
	if payload.Name == "" {
		payload.Name = user.Name
	}
	if payload.Password == "" {
		payload.Password = user.Password
	} else {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("[ERROR] convert string to hash password:", err.Error())
			panic(err)
		}
		payload.Password = string(hashedPassword)
	}
	if payload.Number == nil {
		payload.Number = user.Number
	}
	if payload.Status == nil {
		payload.Status = &user.Status
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	filter := bson.D{
		{Key: "_id", Value: user.ID},
	}
	update := bson.D{
		{Key: "$set", Value: payload},
	}
	coll := r.dm.DB.Collection(UserCollectionName)
	result := coll.FindOneAndUpdate(r.dm.Ctx, filter, update, opts)
	if err := result.Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return user, err
		}
		log.Println("[ERROR] Update document:", err.Error())
		panic(err)
	}

	return user, nil
}

func (r userRepository) UpdateLoginTime(id string) (model.User, error) {
	user, err := r.FindById(id)
	if err != nil {
		return user, err
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	filter := bson.D{
		{Key: "_id", Value: user.ID},
	}
	update := bson.D{
		{Key: "$set", Value: bson.M{
			"lastLoginAt": time.Now().UTC(),
		}},
	}
	coll := r.dm.DB.Collection(UserCollectionName)
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

func (r userRepository) DeleteById(id string) (*mongo.DeleteResult, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, invalidIdErr
	}
	query := bson.D{
		{Key: "_id", Value: _id},
	}

	coll := r.dm.DB.Collection(UserCollectionName)
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
