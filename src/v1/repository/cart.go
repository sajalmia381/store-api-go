package repository

import (
	"log"
	"time"

	"github.com/sajalmia381/store-api/src/enums"
	"github.com/sajalmia381/store-api/src/v1/db"
	"github.com/sajalmia381/store-api/src/v1/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const cartCollectionName = string(enums.CART_COLLECTION_NAME)

type CartRepository interface {
	FindAll() ([]model.Cart, error)
	DeleteByUserId(userId primitive.ObjectID) (*mongo.DeleteResult, error)
	FindByUserId(userId primitive.ObjectID) (model.Cart, error)

	// Requester Cart
	UpdateCartByProduct(userId primitive.ObjectID, payload model.CartProductSpec) (model.Cart, error)
	RemoveProductFromCart(userId primitive.ObjectID, productId primitive.ObjectID) (model.Cart, error)
	UpdateCartByProducts(userId primitive.ObjectID, payload []model.CartProductSpec) (model.Cart, error)
}

type cartRepository struct {
	dm *db.DmManager
}

// Requester Cart
func (r cartRepository) UpdateCartByProducts(userId primitive.ObjectID, products []model.CartProductSpec) (model.Cart, error) {
	var cart model.Cart
	filter := bson.D{
		{Key: "userId", Value: userId},
	}
	payload := bson.D{
		{Key: "$set", Value: bson.M{
			"updatedAt": time.Now().UTC(),
			"products":  products,
		}},
		{Key: "$setOnInsert", Value: bson.M{
			"createdAt": time.Now().UTC(),
		}},
	}
	coll := r.dm.DB.Collection(cartCollectionName)
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	result := coll.FindOneAndUpdate(r.dm.Ctx, filter, payload, opts)
	if err := result.Decode(&cart); err != nil {
		log.Println("[ERROR] product add to cart:", err)
		return cart, err
	}
	return cart, nil
}

func (r cartRepository) UpdateCartByProduct(userId primitive.ObjectID, payload model.CartProductSpec) (model.Cart, error) {
	var cart model.Cart
	filter := bson.D{
		{Key: "userId", Value: userId},
		{Key: "products.productId", Value: payload.ProductId},
	}
	// arrFilter := options.ArrayFilters{
	// 	Filters: bson.A{
	// 		bson.M{"x.productId": payload.ProductId}, // if For Nest Array := bson.M{"y.property_name": filterValue},
	// 	},
	// }
	update := bson.D{
		{Key: "$set", Value: bson.M{
			"updatedAt":  time.Now().UTC(),
			"products.$": payload,
		}},
	}

	upsert := true
	returnDoc := options.After

	opts := &options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &returnDoc,
		// ArrayFilters:   &arrFilter,
	}

	coll := r.dm.DB.Collection(cartCollectionName)
	result := coll.FindOneAndUpdate(r.dm.Ctx, filter, update, opts)
	if err := result.Decode(&cart); err != nil {
		// User Cart not exists
		_filter := bson.D{
			{Key: "userId", Value: userId},
		}
		_update := bson.D{
			{Key: "$push", Value: bson.M{
				"products": payload,
			}},
			{Key: "$set", Value: bson.M{
				"updatedAt": time.Now().UTC(),
			}},
			{
				Key: "$setOnInsert", Value: bson.M{
					"createdAt": time.Now().UTC(),
				},
			},
		}
		_result := coll.FindOneAndUpdate(r.dm.Ctx, _filter, _update, opts)
		if err := _result.Decode(&cart); err != nil {
			log.Println("inner err", err)
			return cart, err
		}
		log.Println("[INFO] User Cart Created")
	}
	return cart, nil
}

func (r cartRepository) RemoveProductFromCart(userId primitive.ObjectID, productId primitive.ObjectID) (model.Cart, error) {
	var cart model.Cart
	filter := bson.D{
		{Key: "userId", Value: userId},
	}
	update := bson.D{
		{Key: "$pull", Value: bson.M{
			"products": bson.M{
				"productId": productId,
			},
		}},
	}
	returnDoc := options.After

	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &returnDoc,
	}

	coll := r.dm.DB.Collection(cartCollectionName)
	result := coll.FindOneAndUpdate(r.dm.Ctx, filter, update, opts)
	if err := result.Decode(&cart); err != nil {
		log.Println("[ERROR] cart remove product: ", err)
		return cart, err
	}
	return cart, nil
}

// Cart CRUD
func (r cartRepository) FindAll() ([]model.Cart, error) {
	var carts []model.Cart

	filter := bson.D{}

	coll := r.dm.DB.Collection(cartCollectionName)
	cursor, err := coll.Find(r.dm.Ctx, filter)
	if err != nil {
		return carts, err
	}
	if err := cursor.All(r.dm.Ctx, &carts); err != nil {
		log.Println("[ERROR] Cart Decade: ", err)
		panic(err)
	}
	return carts, nil
}

func (r cartRepository) FindByUserId(userId primitive.ObjectID) (model.Cart, error) {
	var cart model.Cart
	filter := bson.D{
		{Key: "userId", Value: userId},
	}
	coll := r.dm.DB.Collection(cartCollectionName)
	if err := coll.FindOne(r.dm.Ctx, filter).Decode(&cart); err != nil {
		return cart, err
	}
	return cart, nil
}

func (r cartRepository) DeleteByUserId(userId primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.D{
		{Key: "userId", Value: userId},
	}
	coll := r.dm.DB.Collection(cartCollectionName)
	res, err := coll.DeleteOne(r.dm.Ctx, filter)
	return res, err
}

func NewCartRepository() CartRepository {
	return &cartRepository{
		dm: db.GetDmManager(),
	}
}
