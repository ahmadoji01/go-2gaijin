package middleware

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gitlab.com/kitalabs/go-2gaijin/channels"
	"gitlab.com/kitalabs/go-2gaijin/config"
	"gitlab.com/kitalabs/go-2gaijin/pkg/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB connection string
// for localhost mongoDB
const connectionString = "mongodb://localhost:27017"

// Database Name
const dbName = "go2gaijin"

// Database instance
var DB *mongo.Database

var Pool *websocket.Pool

var IsProduction = false

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
	CreateIndexWithoutWeights(bson.M{"category_ids": 1}, DB.Collection("products"))
}

func HandlePreflight(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
	return
}

func IDExistsInSlice(slice []primitive.ObjectID, val primitive.ObjectID) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
