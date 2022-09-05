package repository

import (
	"errors"
	"log"
	"math"
	"time"

	"github.com/sajalmia381/store-api/src/api/common"
	"github.com/sajalmia381/store-api/src/enums"
	"github.com/sajalmia381/store-api/src/utils"
	"github.com/sajalmia381/store-api/src/v1/db"
	"github.com/sajalmia381/store-api/src/v1/dtos"
	"github.com/sajalmia381/store-api/src/v1/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductRepository interface {
	Store(product model.Product) (model.Product, error)
	FindAll(queryParams dtos.ProductQueryParams) ([]dtos.ProductResponseDto, common.MetaData, error)
	FindBySlug(slug string) (model.Product, error)
	UpdateBySlug(slug string, payload primitive.M) (model.Product, error)
	DeleteBySlug(slug string) (*mongo.DeleteResult, error)
}

type productRepository struct {
	dm *db.DmManager
}

func (p productRepository) Store(product model.Product) (model.Product, error) {
	if product.CreatedBy == "" {
		product.CreatedBy = "anonymous@gmail.com"
	}
	coll := p.dm.DB.Collection(string(enums.PRODUCT_COLLECTION_NAME))
	product.Slug = utils.GenerateUniqueSlug(product.Title, string(enums.PRODUCT_COLLECTION_NAME))
	_, err := coll.InsertOne(p.dm.Ctx, &product)
	if err != nil {
		return product, err
	}
	return product, nil
}

func (p productRepository) FindAll(queryParams dtos.ProductQueryParams) ([]dtos.ProductResponseDto, common.MetaData, error) {
	var objects []dtos.ProductResponseDto
	var aggPipeline mongo.Pipeline
	if queryParams.Search != "" {
		aggPipeline = append(aggPipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "$or", Value: bson.A{
					bson.M{"title": bson.M{"$regex": primitive.Regex{Pattern: queryParams.Search, Options: "i"}}},
					bson.M{"description": bson.M{"$regex": primitive.Regex{Pattern: queryParams.Search, Options: "i"}}},
				}},
			}},
		})
	}

	if queryParams.Sort == enums.DESCENDING {
		aggPipeline = append(aggPipeline, bson.D{
			{Key: "$sort", Value: bson.M{
				"createdAt": -1,
			}},
		})
	}
	coll := p.dm.DB.Collection(string(enums.PRODUCT_COLLECTION_NAME))
	// Pagination
	var metaData common.MetaData
	if queryParams.Limit != 0 {
		// MetaData Calc
		totalProducts, err := coll.CountDocuments(p.dm.Ctx, bson.M{})
		if err != nil {
			panic(err)
		}
		metaData.TotalElements = uint64(totalProducts)
		metaData.PerPage = queryParams.Limit
		if queryParams.Page <= 0 {
			queryParams.Page = 1
		}
		metaData.CurrentPage = queryParams.Page
		metaData.TotalPages = uint64(math.Round(float64(totalProducts) / float64(queryParams.Limit)))
		if metaData.CurrentPage < metaData.TotalPages {
			_nextPage := metaData.CurrentPage + 1
			metaData.NextPage = &_nextPage
		}
		if metaData.CurrentPage > 1 {
			_prevPage := metaData.CurrentPage - 1
			metaData.PrevPage = &_prevPage
		}
		_skipItems := metaData.PerPage * (metaData.CurrentPage - 1)
		skipStage := bson.D{{Key: "$skip", Value: _skipItems}}
		limitStage := bson.D{{Key: "$limit", Value: metaData.PerPage}}
		aggPipeline = append(aggPipeline, skipStage, limitStage)
	}
	// end pagination

	// ForeignKey Aggregation
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

	aggPipeline = append(aggPipeline, categoryLookup, categoryUnwind, userLookup, userUnwind)

	cursor, err := coll.Aggregate(p.dm.Ctx, aggPipeline)
	if err != nil {
		if queryParams.Limit != 0 {
			return objects, metaData, err
		}
		return objects, metaData, err
	}
	if err = cursor.All(p.dm.Ctx, &objects); err != nil {
		log.Println("[ERROR]", err)
		panic(err)
	}
	log.Println("objects", objects)
	if queryParams.Limit != 0 {
		return objects, metaData, nil
	}
	return objects, metaData, nil
}

func (p productRepository) FindBySlug(slug string) (model.Product, error) {
	var product model.Product

	filter := bson.D{
		{Key: "slug", Value: slug},
	}

	coll := p.dm.DB.Collection(string(enums.PRODUCT_COLLECTION_NAME))
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

func (p productRepository) UpdateBySlug(slug string, payload primitive.M) (model.Product, error) {

	filter := bson.D{
		{Key: "slug", Value: slug},
	}
	coll := p.dm.DB.Collection(string(enums.PRODUCT_COLLECTION_NAME))
	payload["updatedAt"] = time.Now().UTC()

	update := bson.M{
		"$set": payload,
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	result := coll.FindOneAndUpdate(p.dm.Ctx, filter, update, opts)
	var product model.Product
	if err := result.Decode(&product); err != nil {
		if err == mongo.ErrNoDocuments {
			return product, errors.New("product is not found")
		}
		log.Println("[ERROR] product update doc:", err)
		panic(err)
	}
	return product, nil
}

func (p productRepository) DeleteBySlug(slug string) (*mongo.DeleteResult, error) {
	filter := bson.D{
		{Key: "slug", Value: slug},
	}
	coll := p.dm.DB.Collection(string(enums.PRODUCT_COLLECTION_NAME))
	result, err := coll.DeleteOne(p.dm.Ctx, filter)
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
