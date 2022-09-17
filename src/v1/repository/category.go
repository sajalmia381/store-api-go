package repository

import (
	"errors"
	"log"

	"github.com/sajalmia381/store-api/src/enums"
	"github.com/sajalmia381/store-api/src/utils"
	"github.com/sajalmia381/store-api/src/v1/db"
	"github.com/sajalmia381/store-api/src/v1/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryRepository interface {
	Store(category model.Category) (model.Category, error)
	FindAll() ([]model.Category, error)
	FindBySlug(slug string) (model.Category, error)
	UpdateBySlug(slug string, payload primitive.M) (model.Category, error)
	DeleteBySlug(slug string) (*mongo.DeleteResult, error)
	PushProductToCategory(categoryId primitive.ObjectID, productId primitive.ObjectID) error
	RemoveProductFromCategory(categoryId primitive.ObjectID, productId primitive.ObjectID) error
	ChangeProductInCategory(oldCategoryId primitive.ObjectID, newCategoryId primitive.ObjectID, productId primitive.ObjectID) error
}

type categoryRepository struct {
	dm *db.DmManager
}

func (r categoryRepository) Store(category model.Category) (model.Category, error) {
	coll := r.dm.DB.Collection(string(enums.CATEGORY_COLLECTION_NAME))
	category.Slug = utils.GenerateUniqueSlug(category.Name, string(enums.CATEGORY_COLLECTION_NAME))
	_, err := coll.InsertOne(r.dm.Ctx, &category)
	if err != nil {
		log.Println("[ERROR] Category Store err: ", err)
		return category, err
	}
	return category, nil
}

func (r categoryRepository) FindAll() ([]model.Category, error) {
	var objects []model.Category
	coll := r.dm.DB.Collection(string(enums.CATEGORY_COLLECTION_NAME))
	filter := bson.D{}
	result, err := coll.Find(r.dm.Ctx, filter)
	if err != nil {
		return objects, err
	}
	if err := result.All(r.dm.Ctx, &objects); err != nil {
		log.Println("[ERROR] category collection cursor", err.Error())
		panic(err)
	}
	return objects, nil
}

func (r categoryRepository) FindBySlug(slug string) (model.Category, error) {
	var category model.Category
	filter := bson.D{
		{Key: "slug", Value: slug},
	}
	coll := r.dm.DB.Collection(string(enums.CATEGORY_COLLECTION_NAME))
	result := coll.FindOne(r.dm.Ctx, filter)
	if err := result.Decode(&category); err != nil {
		if err == mongo.ErrNoDocuments {
			return category, errors.New("category is not found")
		}
		log.Println("[ERROR]:", err.Error())
		panic(err)
	}
	return category, nil
}

func (r categoryRepository) UpdateBySlug(slug string, payload primitive.M) (model.Category, error) {
	var category model.Category
	coll := r.dm.DB.Collection(string(enums.CATEGORY_COLLECTION_NAME))

	filter := bson.D{
		{Key: "slug", Value: slug},
	}
	update := bson.D{
		{Key: "$set", Value: payload},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	result := coll.FindOneAndUpdate(r.dm.Ctx, filter, update, opts)
	if err := result.Decode(&category); err != nil {
		if err == mongo.ErrNoDocuments {
			return category, errors.New("category is not exists")
		}
		log.Println("[ERROR] Update document count", err)
		panic(err)
	}
	return category, nil
}

func (r categoryRepository) DeleteBySlug(slug string) (*mongo.DeleteResult, error) {
	filter := bson.D{
		{Key: "slug", Value: slug},
	}
	coll := r.dm.DB.Collection(string(enums.CATEGORY_COLLECTION_NAME))
	result, err := coll.DeleteOne(r.dm.Ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, errors.New("category is not found! Maybe already deleted")
		}
	}
	if result.DeletedCount != 1 {
		log.Println("[ERROR] Delete document count", result.DeletedCount)
		return result, errors.New("category is not exists")
	}
	return result, nil
}

// CRUD END
func (r categoryRepository) PushProductToCategory(categoryId primitive.ObjectID, productId primitive.ObjectID) error {
	filter := bson.D{
		{Key: "_id", Value: categoryId},
	}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "products", Value: productId},
		}},
	}
	coll := r.dm.DB.Collection(string(enums.CATEGORY_COLLECTION_NAME))
	res, err := coll.UpdateOne(r.dm.Ctx, filter, update)
	if res.MatchedCount != 1 {
		log.Println("Push product to category res:", res)
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("category is not found")
		}
		panic(err)
	}
	return nil
}

func (r categoryRepository) RemoveProductFromCategory(categoryId primitive.ObjectID, productId primitive.ObjectID) error {
	filter := bson.D{
		{Key: "_id", Value: categoryId},
	}
	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "products", Value: productId},
		}},
	}
	coll := r.dm.DB.Collection(string(enums.CATEGORY_COLLECTION_NAME))
	_, err := coll.UpdateOne(r.dm.Ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("category is not found")
		}
		panic(err)
	}
	return nil
}

func (r categoryRepository) ChangeProductInCategory(oldCategoryId primitive.ObjectID, newCategoryId primitive.ObjectID, productId primitive.ObjectID) error {
	oldFilter := bson.D{
		{Key: "_id", Value: oldCategoryId},
	}
	newFilter := bson.D{
		{Key: "_id", Value: newCategoryId},
	}
	oldUpdate := bson.D{
		{Key: "$pull", Value: bson.M{
			"products": productId,
		}},
	}
	newUpdate := bson.D{
		{
			Key: "$push", Value: bson.M{
				"products": productId,
			},
		},
	}
	coll := r.dm.DB.Collection(string(enums.CATEGORY_COLLECTION_NAME))
	_, err := coll.UpdateOne(r.dm.Ctx, oldFilter, oldUpdate)
	_, newErr := coll.UpdateOne(r.dm.Ctx, newFilter, newUpdate)

	if err != nil {
		log.Println("[ERROR] removing product:", err.Error())
	}
	if newErr != nil {
		return err
	}
	return nil
}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{
		dm: db.GetDmManager(),
	}
}
