package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sajalmia381/store-api/src/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var RunMode string
var ServerPort string

var Database string
var DatabaseName string
var DBConnectionString string
var MongoUsername string
var MongoPassword string
var MongoServer string
var MongoPort string

var JwtRegularSecretKey string
var JwtRefreshSecretKey string
var RegularTokenLifetime string

// Super Admin
var SuperAdminName string
var SuperAdminEmail string
var SuperAdminPassword string
var SuperAdminNumber string

// Default User
var DefaultUserId *primitive.ObjectID
var DefaultUserName string
var DefaultUserEmail string
var DefaultUserNumber string
var DefaultUserPassword string

func IntVariables() {
	RunMode = os.Getenv("ENVIRONMENT")
	if RunMode == "" {
		RunMode = string(enums.DEVELOPMENT)
	}

	if RunMode != string(enums.PRODUCTION) {
		err := godotenv.Load()
		if err != nil {
			log.Println("[ERROR]: ", err.Error())
			return
		}
	}
	ServerPort = os.Getenv("SERVER_PORT")
	// Database
	Database = os.Getenv("DATABASE")
	DatabaseName = os.Getenv("DATABASE_NAME")
	MongoServer = os.Getenv("MONGO_SERVER")
	MongoPort = os.Getenv("MONGO_PORT")
	MongoUsername = os.Getenv("MONGO_USERNAME")
	MongoPassword = os.Getenv("MONGO_PASSWORD")

	if Database == "MONGO" {
		DBConnectionString = "mongodb://" + MongoUsername + ":" + MongoPassword + "@" + MongoServer + ":" + MongoPort
	}

	fmt.Printf("DB Username: %s, DB Server: %s, DB Port: %s", MongoUsername, MongoServer, MongoPort)

	// JWT
	JwtRegularSecretKey = os.Getenv("JWT_SECRET_KEY")
	JwtRefreshSecretKey = os.Getenv("JWT_REFRESH_KEY")
	RegularTokenLifetime = os.Getenv("REGULAR_TOKEN_LIFETIME")

	// Super Admin
	SuperAdminName = os.Getenv("SUPER_ADMIN_NAME")
	SuperAdminEmail = os.Getenv("SUPER_ADMIN_EMAIL")
	SuperAdminPassword = os.Getenv("SUPER_ADMIN_PASSWORD")
	SuperAdminNumber = os.Getenv("SUPER_ADMIN_NUMBER")

	// Default User
	DefaultUserName = "Anonymous User"
	DefaultUserEmail = "anonymous@gmail.com"
	DefaultUserNumber = "1234567891"
	DefaultUserPassword = "simple_password"
}
