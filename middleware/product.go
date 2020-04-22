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
	payload := getHome()
	json.NewEncoder(c.Writer).Encode(payload)
}

func GetProductDetail(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	var collection = DB.Collection("products")

	productID, err := primitive.ObjectIDFromHex(c.Param("id"))

	var res models.ResponseResult
	if err != nil {
		res.Error = "Error while converting ID, try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	var payload models.Product
	err = collection.FindOne(context.Background(), bson.M{"_id": productID}).Decode(&payload)
	if err != nil {
		res.Error = "Error while searching for product, try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	json.NewEncoder(c.Writer).Encode(payload)
}

func getHome() responses.HomePage {
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
		collection = DB.Collection("categories")
		homeData.Categories = PopulateCategories(collection.Find(context.Background(), options))
		wg.Done()
	}()

	// Get Recent Items
	wg.Add(1)
	go func() {
		collection = DB.Collection("products")
		options.SetLimit(16)
		options.SetSort(bson.D{{"created_at", -1}})
		homeData.RecentItems = PopulateProducts(collection.Find(context.Background(), bson.D{{"status_cd", 1}}, options))
		wg.Done()
	}()

	// Get Free Items
	wg.Add(1)
	go func() {
		collection = DB.Collection("products")
		options.SetLimit(16)
		options.SetSort(bson.D{{"created_at", -1}})
		homeData.FreeItems = PopulateProducts(collection.Find(context.Background(), bson.D{{"price", 0}}, options))
		wg.Done()
	}()

	// Get Recommended Items
	wg.Add(1)
	go func() {
		collection = DB.Collection("products")
		options.SetLimit(16)
		options.SetSort(bson.D{{"created_at", -1}})
		homeData.RecommendedItems = PopulateProducts(collection.Find(context.Background(), bson.D{{"price", 0}}, options))
		wg.Done()
	}()

	// Get Featured Items
	wg.Add(1)
	go func() {
		collection = DB.Collection("products")
		options.SetLimit(16)
		options.SetSort(bson.D{{"created_at", -1}})
		homeData.FeaturedItems = PopulateProducts(collection.Find(context.Background(), bson.D{{"price", 0}}, options))
		wg.Done()
	}()

	wg.Wait()

	homeResponse.Data = homeData
	homeResponse.Status = "Success"
	homeResponse.Message = "Homepage Data Loaded"

	return homeResponse
}
