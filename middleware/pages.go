package middleware

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/gin-gonic/gin"
	"gitlab.com/kitalabs/go-2gaijin/models"
	"gitlab.com/kitalabs/go-2gaijin/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetHome(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	var wg sync.WaitGroup
	var homeResponse responses.HomePage
	var homeData responses.HomeData
	var collection = DB.Collection("")

	var options = &options.FindOptions{}

	// Get Banners
	wg.Add(1)
	go func() {
		collection = DB.Collection("banners")
		options.SetLimit(5)
		homeData.Banners = PopulateBanners(collection.Find(context.Background(), bson.D{{}}, options))
		wg.Done()
	}()

	// Get Categories
	wg.Add(1)
	go func() {
		locale := "en"
		homeData.Categories = PopulateCategories(locale)
		wg.Done()
	}()

	// Get Recent Items
	wg.Add(1)
	go func() {
		projection := bson.D{{"_id", 1}, {"name", 1}, {"price", 1}, {"img_url", 1}, {"user_id", 1}, {"seller_name", 1}, {"location", 1}, {"status_cd", 1}}
		sort := bson.D{{"created_at", -1}}

		options.SetProjection(projection)
		options.SetSort(sort)
		options.SetLimit(16)

		homeData.RecentItems = PopulateProductsWithAnImage(bson.D{}, options)
		wg.Done()
	}()

	// Get Free Items
	wg.Add(1)
	go func() {
		projection := bson.D{{"_id", 1}, {"name", 1}, {"price", 1}, {"img_url", 1}, {"user_id", 1}, {"seller_name", 1}, {"location", 1}, {"status_cd", 1}}
		sort := bson.D{{"created_at", -1}}

		options.SetProjection(projection)
		options.SetSort(sort)
		options.SetLimit(16)

		homeData.FreeItems = PopulateProductsWithAnImage(bson.D{{"price", 0}}, options)
		wg.Done()
	}()

	// Get Recommended Items
	wg.Add(1)
	go func() {
		projection := bson.D{{"_id", 1}, {"name", 1}, {"price", 1}, {"img_url", 1}, {"user_id", 1}, {"seller_name", 1}, {"location", 1}, {"status_cd", 1}}
		sort := bson.D{{"created_at", -1}}

		options.SetProjection(projection)
		options.SetSort(sort)
		options.SetLimit(16)

		homeData.RecommendedItems = PopulateProductsWithAnImage(bson.D{{"price", 0}}, options)
		wg.Done()
	}()

	// Get Featured Items
	wg.Add(1)
	go func() {
		projection := bson.D{{"_id", 1}, {"name", 1}, {"price", 1}, {"img_url", 1}, {"user_id", 1}, {"seller_name", 1}, {"location", 1}, {"status_cd", 1}}
		sort := bson.D{{"created_at", -1}}

		options.SetProjection(projection)
		options.SetSort(sort)
		options.SetLimit(16)

		homeData.FeaturedItems = PopulateProductsWithAnImage(bson.D{{"price", 0}}, options)
		wg.Done()
	}()

	wg.Wait()

	homeResponse.Data = homeData
	homeResponse.Status = "Success"
	homeResponse.Message = "Homepage Data Loaded"

	json.NewEncoder(c.Writer).Encode(homeResponse)
}

func GetProductDetail(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	var collection = DB.Collection("products")
	var locale = "en"

	productID, err := primitive.ObjectIDFromHex(c.Param("id"))

	var res models.ResponseResult
	if err != nil {
		res.Error = "Error while converting ID, try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	var payload responses.ProductDetailPage
	err = collection.FindOne(context.Background(), bson.M{"_id": productID}).Decode(&payload)
	payload.Category = FindACategoryFromProductID(productID, locale)
	payload.Status = ProductStatusEnum(payload.StatusEnum)

	if err != nil {
		res.Error = "Error while searching for product, try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	json.NewEncoder(c.Writer).Encode(payload)
}

func GetChatLobby(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	var roomsData []models.Room

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		roomsData = PopulateRoomsFromUserID(userData.ID)
		if roomsData != nil {
			json.NewEncoder(c.Writer).Encode(roomsData)
		}
		return
	}
}
