package middleware

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gitlab.com/kitalabs/go-2gaijin/pkg/websocket"
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

// create connection with mongo DB
func init() {

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	DB = client.Database(dbName)

	productIndex()

}

func productIndex() {
	weights := bson.M{"name": 5, "description": 2}
	keys := bson.M{"name": "text", "description": "text"}
	CreateIndex(weights, keys, DB.Collection("products"))
}

func serveWs(pool *websocket.Pool, c *gin.Context) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(c.Writer, c.Request)
	if err != nil {
		fmt.Fprintf(c.Writer, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func WebSocketHandler(c *gin.Context) {
	pool := websocket.NewPool()
	go pool.Start()

	serveWs(pool, c)
}
