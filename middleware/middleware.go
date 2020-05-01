package middleware

import (
	"context"
	"fmt"
	"log"

	"gitlab.com/kitalabs/go-2gaijin/channels"
	"gitlab.com/kitalabs/go-2gaijin/pkg/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB connection string
// for localhost mongoDB
const connectionString = "mongodb://localhost:27017"

// Database Name
const dbName = "go-2gaijin"

// Database instance
var DB *mongo.Database

//Image URL Prefix
var ImgURLPrefix string = "https://storage.googleapis.com/rails-2gaijin-storage/"

var Pool *websocket.Pool

// create connection with mongo DB
func init() {

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	DB = client.Database(dbName)

	// Index product
	productIndex()

	go channels.H.Run()
}

func productIndex() {
	weights := bson.M{"name": 5, "description": 2}
	keys := bson.M{"name": "text", "description": "text"}
	CreateIndex(weights, keys, DB.Collection("products"))
}
