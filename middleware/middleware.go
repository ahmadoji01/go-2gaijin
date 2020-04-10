package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	jwt "github.com/dgrijalva/jwt-go"
	"gitlab.com/kitalabs/go-2gaijin/models"
	"gitlab.com/kitalabs/go-2gaijin/pkg/websocket"
	"gitlab.com/kitalabs/go-2gaijin/responses"
	"golang.org/x/crypto/bcrypt"
)

// DB connection string
// for localhost mongoDB
const connectionString = "mongodb://localhost:27017"

// Database Name
const dbName = "go-2gaijin"

// Collection name
const collName = "products"

// collection object/instance
var collection *mongo.Collection

// Database instance
var db *mongo.Database

// create connection with mongo db
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

	db = client.Database(dbName)

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

func GetSearch(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	urlQuery := c.Request.URL.Query()

	query := urlQuery.Get("q")
	sort := urlQuery.Get("sortby")
	asc, err := strconv.Atoi(urlQuery.Get("asc"))
	if sort == "" {
		sort = "created_at"
		asc = -1
	}

	var status = -1
	strStatus := urlQuery.Get("status")
	if strStatus == "sold" {
		status = 1
	} else if strStatus == "available" {
		status = 0
	} else {
		status = -1
	}

	category := urlQuery.Get("category")
	//asc := c.Param("asc")
	//status := c.Param("status")
	var limit int64
	var page int64
	if urlQuery.Get("limit") == "" {
		limit = 16
	} else {
		limit, err = strconv.ParseInt(urlQuery.Get("limit"), 10, 64)
	}

	if urlQuery.Get("page") == "" {
		page = 1
	} else {
		page, err = strconv.ParseInt(urlQuery.Get("page"), 10, 64)
	}

	var res models.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	payload := getSearch(query, category, limit, page, sort, asc, status)
	json.NewEncoder(c.Writer).Encode(payload)
}

func GetHome(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	payload := getHome()
	json.NewEncoder(c.Writer).Encode(payload)
}

// GetAllTask get all the task route
func GetAllProducts(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	payload := getAllProducts()
	json.NewEncoder(c.Writer).Encode(payload)
}

func RegisterHandler(c *gin.Context) {

	c.Writer.Header().Set("Content-Type", "application/json")
	var user models.User
	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &user)
	var res models.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	collection := db.Collection("users")

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	var result models.User
	err = collection.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Decode(&result)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

			if err != nil {
				res.Error = "Error While Hashing Password, Try Again"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}
			user.Password = string(hash)

			_, err = collection.InsertOne(context.TODO(), user)
			if err != nil {
				res.Error = "Error While Creating User, Try Again"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}
			res.Result = "Registration Successful"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		res.Error = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	res.Result = "Email already Exists!!"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func LoginHandler(c *gin.Context) {

	c.Writer.Header().Set("Content-Type", "application/json")
	var user models.User
	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}

	collection := db.Collection("users")

	if err != nil {
		log.Fatal(err)
	}
	var result models.User
	var res models.ResponseResult

	err = collection.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Decode(&result)

	if err != nil {
		res.Error = "Invalid email"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))

	if err != nil {
		res.Error = "Invalid password"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":      result.Email,
		"first_name": result.FirstName,
		"last_name":  result.LastName,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("MY_JWT_TOKEN")))

	if err != nil {
		res.Error = "Error while generating token, try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	result.Token = tokenString
	result.Password = ""

	json.NewEncoder(c.Writer).Encode(result)

}

func ProfileHandler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(os.Getenv("MY_JWT_TOKEN")), nil
	})
	var result models.User
	var res models.ResponseResult
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Email = claims["email"].(string)
		result.FirstName = claims["first_name"].(string)
		result.LastName = claims["last_name"].(string)

		json.NewEncoder(c.Writer).Encode(result)
		return
	} else {
		res.Error = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

}

func getSearch(query string, category string, nPerPage int64, page int64, sort string, asc int, status int) responses.SearchPage {

	var wg sync.WaitGroup
	var filter bson.M
	var collection = db.Collection("products")
	fmt.Println(query)

	filter = searchFilter(query, status)

	var options = &options.FindOptions{}
	options.SetLimit(nPerPage)
	options.SetSort(bson.D{{sort, asc}})
	options.SetSkip(nPerPage * (page - 1))

	var pagination responses.Pagination
	var searchData responses.SearchData
	var searchResponse responses.SearchPage
	var items []models.Product

	wg.Add(1)
	go func() {
		items = populateProducts(collection.Find(context.Background(), filter, options))
		wg.Done()
	}()

	var count int64
	var err error
	wg.Add(1)
	go func() {
		count, err = collection.CountDocuments(context.Background(), filter)
		wg.Done()
	}()
	wg.Wait()

	if err != nil {
		searchResponse.Status = "Error"
		searchResponse.Message = "Error Counting Documents. Try Again"
		return searchResponse
	}

	pagination.CurrentPage = page
	pagination.ItemsPerPage = nPerPage
	pagination.TotalItems = count

	if count-(nPerPage*(page+1)) >= 1 {
		pagination.NextPage = page + 1
	} else {
		pagination.NextPage = 0
	}

	if page-1 >= 1 {
		pagination.PreviousPage = page - 1
	} else {
		pagination.PreviousPage = 0
	}

	searchData.Items = items
	searchData.Pagination = pagination

	searchResponse.Data = searchData
	searchResponse.Status = "Success"
	searchResponse.Message = "Search Page Data Loaded"

	return searchResponse
}

func searchFilter(query string, status int) bson.M {
	var filter bson.M

	if query != "" && status != -1 {
		filter = bson.M{"name": query, "status_cd": status}
	} else if query == "" && status != -1 {
		filter = bson.M{"status_cd": status}
	} else if query != "" && status == -1 {
		filter = bson.M{"name": query}
	} else {
		filter = bson.M{}
	}

	return filter

}

func getHome() responses.HomePage {
	var wg sync.WaitGroup
	var homeResponse responses.HomePage
	var homeData responses.HomeData
	var collection = db.Collection("")

	var options = &options.FindOptions{}

	// Get Banners
	wg.Add(1)
	go func() {
		collection = db.Collection("banners")
		options.SetLimit(5)
		homeData.Banners = populateBanners(collection.Find(context.Background(), bson.D{{}}, options))
		wg.Done()
	}()

	// Get Categories
	wg.Add(1)
	go func() {
		collection = db.Collection("categories")
		homeData.Categories = populateCategories(collection.Find(context.Background(), bson.D{{}}, options))
		wg.Done()
	}()

	// Get Recent Items
	wg.Add(1)
	go func() {
		collection = db.Collection("products")
		options.SetLimit(16)
		options.SetSort(bson.D{{"created_at", -1}})
		homeData.RecentItems = populateProducts(collection.Find(context.Background(), bson.D{{"status_cd", 1}}, options))
		wg.Done()
	}()

	// Get Free Items
	wg.Add(1)
	go func() {
		collection = db.Collection("products")
		options.SetLimit(16)
		options.SetSort(bson.D{{"created_at", -1}})
		homeData.FreeItems = populateProducts(collection.Find(context.Background(), bson.D{{"price", 0}}, options))
		wg.Done()
	}()

	// Get Recommended Items
	wg.Add(1)
	go func() {
		collection = db.Collection("products")
		options.SetLimit(16)
		options.SetSort(bson.D{{"created_at", -1}})
		homeData.RecommendedItems = populateProducts(collection.Find(context.Background(), bson.D{{"price", 0}}, options))
		wg.Done()
	}()

	// Get Featured Items
	wg.Add(1)
	go func() {
		collection = db.Collection("products")
		options.SetLimit(16)
		options.SetSort(bson.D{{"created_at", -1}})
		homeData.FeaturedItems = populateProducts(collection.Find(context.Background(), bson.D{{"price", 0}}, options))
		wg.Done()
	}()

	wg.Wait()

	homeResponse.Data = homeData
	homeResponse.Status = "Success"
	homeResponse.Message = "Homepage Data Loaded"

	return homeResponse
}

func populateProducts(cur *mongo.Cursor, err error) []models.Product {
	var results []models.Product
	for cur.Next(context.Background()) {
		var result models.Product
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

func populateBanners(cur *mongo.Cursor, err error) []models.Banner {
	var results []models.Banner
	for cur.Next(context.Background()) {
		var result models.Banner
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

func populateCategories(cur *mongo.Cursor, err error) []models.Category {
	var results []models.Category
	for cur.Next(context.Background()) {
		var result models.Category
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

// get all task from the DB and return it
func getAllProducts() []primitive.M {
	var collection = db.Collection("products")
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}