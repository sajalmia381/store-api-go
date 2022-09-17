package db

import (
	"context"
	"log"
	"sync"

	"github.com/sajalmia381/store-api/src/config"
	"github.com/sajalmia381/store-api/src/enums"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DmManager struct {
	Ctx context.Context
	DB  *mongo.Database
}

var singletonDmManager *DmManager
var onceDmManager sync.Once

func GetDmManager() *DmManager {
	onceDmManager.Do(func() {
		singletonDmManager = &DmManager{}
		singletonDmManager.initializeConnection()
	})
	return singletonDmManager
}

func (dm *DmManager) initializeConnection() {
	ctx := context.Background()
	dm.Ctx = ctx
	clientOpts := options.Client().ApplyURI(config.DBConnectionString)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Println("[ERROR] SingletonDB connection error: ", err.Error())
		return
	}
	db := client.Database(config.DatabaseName)
	collectionNames, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		log.Println("GETTING collection name err:", collectionNames)
	}
	for _, _name := range enums.COLLECTION_NAMES {
		isExists := false
		for _, collName := range collectionNames {
			if collName == _name {
				isExists = true
				break
			}
		}
		if !isExists {
			db.CreateCollection(ctx, _name)
			switch _name {
			case string(enums.USER_COLLECTION_NAME):
				{
					indexModel := mongo.IndexModel{
						Keys:    bson.D{{Key: "email", Value: 1}},
						Options: options.Index().SetUnique(true),
					}
					// indexModelNumber := mongo.IndexModel{
					// 	Keys:    bson.D{{Key: "number", Value: 1}},
					// 	Options: options.Index().SetUnique(true),
					// }
					_names, _ := db.Collection(_name).Indexes().CreateMany(context.Background(), []mongo.IndexModel{
						indexModel, // indexModelNumber,
					})
					log.Println("Index created! indexNames: ", _names)
				}
			case string(enums.CART_COLLECTION_NAME):
				{
					indexModel := mongo.IndexModel{
						Keys:    bson.D{{Key: "userId", Value: 1}},
						Options: options.Index().SetUnique(true),
					}
					_name, _ := db.Collection(_name).Indexes().CreateOne(context.Background(), indexModel)
					log.Println("Index created! indexName: ", _name)
				}
			case string(enums.PRODUCT_COLLECTION_NAME):
				{
					indexModel := mongo.IndexModel{
						Keys:    bson.D{{Key: "slug", Value: 1}},
						Options: options.Index().SetUnique(true),
					}
					_name, _ := db.Collection(_name).Indexes().CreateOne(context.Background(), indexModel)
					log.Println("Index created! indexName: ", _name)
				}
			case string(enums.CATEGORY_COLLECTION_NAME):
				{
					indexModel := mongo.IndexModel{
						Keys:    bson.D{{Key: "slug", Value: 1}},
						Options: options.Index().SetUnique(true),
					}
					_name, _ := db.Collection(_name).Indexes().CreateOne(context.Background(), indexModel)
					log.Println("Index created! indexName: ", _name)
				}
			case string(enums.TOKEN_COLLECTION_NAME):
				{
					indexModel := mongo.IndexModel{
						Keys:    bson.D{{Key: "token", Value: 1}},
						Options: options.Index().SetUnique(true),
					}
					_name, _ := db.Collection(_name).Indexes().CreateOne(context.Background(), indexModel)

					log.Println("Index created! indexName: ", _name)
				}
			}

		}
	}
	dm.DB = db
	log.Println("[INFO] Initialized Singleton DB Manager")
}
