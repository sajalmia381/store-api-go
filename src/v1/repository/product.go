package repository

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/sajalmia381/store-api/src/utils"
	"github.com/sajalmia381/store-api/src/v1/db"
	"github.com/sajalmia381/store-api/src/v1/dtos"
	"github.com/sajalmia381/store-api/src/v1/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ProductCollectionName = "products"
)

type ProductRepository interface {
	Store(product model.Product) (model.Product, error)
	FindAll() ([]dtos.ProductResponseDto, error)
	FindBySlug(slug string) (model.Product, error)
	UpdateOne(product model.Product, payload dtos.ProductUpdateDto) (model.Product, error)
	DeleteBySlug(slug string) (*mongo.DeleteResult, error)
}

type productRepository struct {
	dm *db.DmManager
}

func (p productRepository) Store(product model.Product) (model.Product, error) {
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now().UTC()
	product.UpdatedAt = time.Now().UTC()
	product.Active = true
	if product.CreatedBy == "" {
		product.CreatedBy = "anonymous@gmail.com"
	}

	coll := p.dm.DB.Collection(ProductCollectionName)
	product.Slug = utils.GenerateUniqueSlug(product.Title, coll)
	_, err := coll.InsertOne(p.dm.Ctx, &product)
	if err != nil {
		return product, err
	}
	return product, nil
}

func (p productRepository) FindAll() ([]dtos.ProductResponseDto, error) {
	var objects []dtos.ProductResponseDto

	categoryLookup := bson.D{
		{
			Key: "$lookup", Value: bson.M{
				"from":         "categories", // the collection name
				"localField":   "category",   // the field on the child struct
				"foreignField": "_id",        // the field on the parent struct
				"as":           "category",   // the field to populate into
			},
		},
	}

	categoryUnwind := bson.D{
		{Key: "$unwind", Value: bson.M{
			"path":                       "$category",
			"preserveNullAndEmptyArrays": true,
		}},
	}
	userLookup := bson.D{
		{
			Key: "$lookup", Value: bson.M{
				"from":         "users",
				"localField":   "createdBy",
				"foreignField": "email",
				"as":           "createdBy",
			},
		},
	}

	userUnwind := bson.D{
		{
			Key: "$unwind", Value: bson.M{
				"path":                       "$createdBy",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}

	aggPipeline := mongo.Pipeline{
		categoryLookup,
		categoryUnwind,
		userLookup,
		userUnwind,
	}

	coll := p.dm.DB.Collection(ProductCollectionName)
	cursor, err := coll.Aggregate(p.dm.Ctx, aggPipeline)
	if err != nil {
		return objects, err
	}

	if err = cursor.All(p.dm.Ctx, &objects); err != nil {
		log.Println("[ERROR]", err)
		panic(err)
	}
	return objects, nil
}

func (p productRepository) FindBySlug(slug string) (model.Product, error) {
	var product model.Product

	filter := bson.D{
		{Key: "slug", Value: slug},
	}

	coll := p.dm.DB.Collection(ProductCollectionName)
	result := coll.FindOne(p.dm.Ctx, filter)
	if err := result.Decode(&product); err != nil {
		if err == mongo.ErrNoDocuments {
			return product, errors.New("product is not found")
		}
		log.Println("[ERROR]", err)
		panic(err)
	}
	return product, nil
}

func (p productRepository) UpdateOne(product model.Product, payload dtos.ProductUpdateDto) (model.Product, error) {

	filter := bson.D{
		{Key: "slug", Value: product.Slug},
	}
	coll := p.dm.DB.Collection(ProductCollectionName)
	payloadMap := map[string]interface{}{
		"updateAt": time.Now().UTC(),
	}
	if payload.Title != "" {
		payloadMap["title"] = payload.Title
	}
	if payload.Price != nil {
		payloadMap["price"] = payload.Price
	}
	if payload.Description != "" {
		payloadMap["description"] = payload.Description
	}
	if payload.Category != nil {
		payloadMap["category"] = payload.Category
	}

	if payload.UpdateSlug {
		if payload.Title != "" {
			payloadMap["slug"] = utils.GenerateUniqueSlug(payload.Title, coll, product.Slug)
		} else {
			payloadMap["slug"] = utils.GenerateUniqueSlug(product.Title, coll, product.Slug)
		}
	}
	update := bson.M{
		"$set": payloadMap,
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	result := coll.FindOneAndUpdate(p.dm.Ctx, filter, update, opts)
	if err := result.Decode(&product); err != nil {
		log.Println("[ERROR] product update doc:", err)
		panic(err)
	}
	return product, nil
}

func (p productRepository) DeleteBySlug(slug string) (*mongo.DeleteResult, error) {
	filter := bson.D{
		{Key: "slug", Value: slug},
	}
	coll := p.dm.DB.Collection(ProductCollectionName)
	result, err := coll.DeleteOne(p.dm.Ctx, filter)
	fmt.Println("result", result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, errors.New("product is not exists")
		}
		panic(err)
	}
	return result, nil
}

func NewProductRepository() ProductRepository {
	return &productRepository{
		dm: db.GetDmManager(),
	}
}
