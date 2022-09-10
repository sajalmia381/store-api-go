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
		for _, dbName := range collectionNames {
			if dbName != _name {
				db.CreateCollection(ctx, _name)
			}
		}
	}
	dm.DB = db
	log.Println("[INFO] Initialized Singleton DB Manager")
}
