package middleware

import (
	"context"
	"encoding/json"
	"strconv"
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
	var wg sync.WaitGroup

	productID, err := primitive.ObjectIDFromHex(c.Param("id"))

	var res responses.ResponseMessage
	if err != nil {
		res.Status = "Error"
		res.Message = "Error while converting ID, try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	var output responses.ProductDetailPage
	var payload responses.ProductDetails
	var item responses.ProductDetailItem
	err = collection.FindOne(context.Background(), bson.M{"_id": productID}).Decode(&item)

	if err != nil {
		res.Status = "Error"
		res.Message = "Error while retrieving product info, try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	var location = struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}{}

	latitude, err := strconv.ParseFloat(item.Latitude, 64)
	longitude, err := strconv.ParseFloat(item.Longitude, 64)

	if err != nil {
		res.Status = "Error"
		res.Message = "Error while converting numbers, try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	location.Latitude = latitude
	location.Longitude = longitude

	if len(item.Comments) == 0 {
		item.Comments = make([]interface{}, 0)
	}
	item.Location = location
	item.Category = FindACategoryFromProductID(productID, locale)
	item.Status = ProductStatusEnum(item.StatusEnum)
	item.Images = FindProductImages(productID)
	payload.Item = item

	var options = &options.FindOptions{}
	projection := bson.D{{"_id", 1}, {"name", 1}, {"price", 1}, {"img_url", 1}, {"user_id", 1}, {"seller_name", 1}, {"latitude", 1}, {"longitude", 1}, {"status_cd", 1}}
	sort := bson.D{{"created_at", -1}}
	options.SetLimit(8)
	options.SetProjection(projection)
	options.SetSort(sort)

	// Search Related Items
	wg.Add(1)
	go func() {
		filter := bson.D{{"category_ids", item.Category.ID}}
		payload.RelatedItems = PopulateProductsWithAnImage(filter, options)
		wg.Done()
	}()

	// Search Seller Items
	wg.Add(1)
	go func() {
		filter := bson.D{{"user_id", item.User}}
		payload.SellerItems = PopulateProductsWithAnImage(filter, options)
		wg.Done()
	}()
	wg.Wait()

	if err != nil {
		res.Status = "Error"
		res.Message = "Error while searching for product, try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	output.Status = "Success"
	output.Message = "Product Detail Successfully Loaded"
	output.Data = payload

	json.NewEncoder(c.Writer).Encode(output)
}

func GetChatLobby(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	var roomsData []models.Room
	var lobbyData responses.ChatLobbyData
	var res responses.GeneralResponse

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		roomsData = PopulateRoomsFromUserID(userData.ID)
		res.Message = "Chat Lobby Retrieved!"
		res.Status = "Success"
		if roomsData == nil {
			roomsData = make([]models.Room, 0)
		}

		lobbyData.ChatLobby = roomsData
		res.Data = lobbyData
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
}

func GetWishlistPage(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	var start int64
	var limit int64
	var err error

	if c.Request.URL.Query().Get("start") == "" {
		start = 0
	} else {
		start, err = strconv.ParseInt(c.Request.URL.Query().Get("start"), 10, 64)
	}
	if c.Request.URL.Query().Get("limit") == "" {
		limit = 8
	} else {
		limit, err = strconv.ParseInt(c.Request.URL.Query().Get("limit"), 10, 64)
	}

	if err != nil {
		var res responses.ResponseMessage
		res.Status = "Error"
		res.Message = "Error converting data. Try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)

	if !isLoggedIn {
		var res responses.ResponseMessage
		res.Status = "Error"
		res.Message = "You have to login to use this feature"
		json.NewEncoder(c.Writer).Encode(res)
		return
	} else {
		filter := bson.D{{"user_id", userData.ID}}
		options := &options.FindOptions{}
		options.SetSkip(start)
		options.SetLimit(limit)

		wishlistItems := PopulateProductsWithAnImage(filter, options)

		var res responses.SearchPage
		res.Status = "Success"
		res.Message = "Wishlist items have been successfully retrieved"
		res.Data = wishlistItems
		if wishlistItems == nil {
			res.Data = make([]interface{}, 0)
		}

		json.NewEncoder(c.Writer).Encode(res)
		return
	}
}