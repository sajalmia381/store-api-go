package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client
var isDatabaseConnected bool

func InitDBConnection() bool {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DBConnectionString))
	if err != nil {
		log.Println("[ERROR] Database Connection Failed: ", err)
		return false
	}
	mongoClient = client
	isDatabaseConnected = true
	return true
}

func CloseConnection() {
	if err := mongoClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func ping() bool {
	if isDatabaseConnected && mongoClient != nil {
		if err := mongoClient.Ping(context.TODO(), readpref.Primary()); err != nil {
			log.Println("[ERROR] Database ping failed: ", err.Error())
			return false
		}
		//log.Println("[INFO] Successfully connected and pinged.")
		return true
	}
	log.Println("[ERROR] Database connection does not exist!")
	return false
}

func reconnect() {
	for !InitDBConnection() {
		log.Println("[INFO] Try to reconnect database...")
		time.Sleep(time.Duration(2) * time.Second)
	}
	isDatabaseConnected = true
}

func DBHealthChecker() {
	log.Println("[INFO] Start database health checker.")
	for {
		if !ping() {
			isDatabaseConnected = false
			break
		}
		time.Sleep(time.Duration(5) * time.Second)
	}
	reconnect()
}
