package middleware

import (
	"context"
	"encoding/json"
	"fmt"
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
		homeData.Categories = PopulateCategories()
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
		userID, _ := primitive.ObjectIDFromHex("5da95727697d19f3f01f62b6")

		projection := bson.D{{"_id", 1}, {"name", 1}, {"price", 1}, {"img_url", 1}, {"user_id", 1}, {"seller_name", 1}, {"location", 1}, {"status_cd", 1}}
		sort := bson.D{{"created_at", -1}}

		options.SetProjection(projection)
		options.SetSort(sort)
		options.SetLimit(4)

		homeData.FeaturedItems = PopulateProductsWithAnImage(bson.D{{"user_id", userID}}, options)
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

	latitude := item.Latitude
	longitude := item.Longitude

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

	wg.Add(1)
	go func() {
		var cat models.Category
		collection = DB.Collection("categories")
		var findOneOpt = &options.FindOneOptions{}
		findOneOpt.SetProjection(bson.D{{"_id", 1}, {"name", 1}, {"icon_url", 1}})
		err = collection.FindOne(context.Background(), bson.M{"_id": item.CategoryIDs[0]}, findOneOpt).Decode(&cat)
		if err != nil {
			res.Status = "Error"
			res.Message = "Error while retrieving product info, try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}
		item.Category = cat
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		var detail models.ProductDetail
		collection = DB.Collection("product_details")
		collection.FindOne(context.Background(), bson.M{"product_id": productID}).Decode(&detail)
		payload.Detail = detail
		wg.Done()
	}()

	wg.Wait()

	item.Location = location
	item.Status = ProductStatusEnum(item.StatusEnum)
	item.Images = FindProductImages(productID)
	payload.Item = item

	var options = &options.FindOptions{}
	projection := bson.D{{"_id", 1}, {"name", 1}, {"price", 1}, {"img_url", 1}, {"user_id", 1}, {"seller_name", 1}, {"latitude", 1}, {"longitude", 1}, {"status_cd", 1}}
	sort := bson.D{{"created_at", -1}}
	options.SetLimit(8)
	options.SetProjection(projection)
	options.SetSort(sort)

	// Get Seller Info
	wg.Add(1)
	go func() {
		payload.Seller = GetSellerInfo(item.User)
		wg.Done()
	}()

	// Search Related Items
	wg.Add(1)
	go func() {
		filter := bson.D{{"category_ids", item.CategoryIDs[0]}}
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
	fmt.Println(output)

	json.NewEncoder(c.Writer).Encode(output)
}

func GetChatLobby(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	var res responses.GenericResponse
	var err error

	urlQuery := c.Request.URL.Query()
	var start int64
	var limit int64
	if urlQuery.Get("start") == "" {
		start = 0
	} else {
		start, err = strconv.ParseInt(urlQuery.Get("start"), 10, 64)
	}

	if urlQuery.Get("limit") == "" {
		limit = 8
	} else {
		limit, err = strconv.ParseInt(urlQuery.Get("limit"), 10, 64)
		if limit <= 0 {
			limit = 8
		}
	}

	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	var roomsData []models.Room
	var lobbyData responses.ChatLobbyData

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		roomsData, _ = PopulateRoomsFromUserID(userData.ID, start, limit)
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
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func GetSellerAppointmentPage(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)

	var res responses.GenericResponse
	var appointments []models.Appointment
	var appointmentData responses.AppointmentData

	if isLoggedIn {
		appointments = PopulateAppointmentsFromUserID(userData.ID, "seller")
		appointmentData.Appointments = appointments

		res.Status = "Success"
		res.Message = "Appointments Retrieved"
		res.Data = appointmentData
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func GetBuyerAppointmentPage(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)

	var res responses.GenericResponse
	var appointments []models.Appointment
	var appointmentData responses.AppointmentData

	if isLoggedIn {
		appointments = PopulateAppointmentsFromUserID(userData.ID, "buyer")
		appointmentData.Appointments = appointments

		res.Status = "Success"
		res.Message = "Appointments Retrieved"
		res.Data = appointmentData
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func GetNotificationPage(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)

	var res responses.GenericResponse
	var notifications []models.Notification
	var notificationData responses.NotificationData

	if isLoggedIn {
		idFilter := bson.D{{"notified_id", userData.ID}}
		notifications = PopulateNotificationsFromUserID(idFilter)
		notificationData.Notifications = notifications

		res.Status = "Success"
		res.Message = "Notifications Retrieved"
		res.Data = notificationData
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
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

func GetProfileForVisitorPage(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	var res responses.GenericResponse
	var userID = c.Request.URL.Query().Get("user_id")
	var wg sync.WaitGroup

	if userID == "" {
		res.Status = "Error"
		res.Message = "You have to input the user ID"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		res.Status = "Error"
		res.Message = "Something went wrong"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	var profileData responses.ProfileForVisitorData
	var opts = &options.FindOptions{}
	projection := bson.D{{"_id", 1}, {"name", 1}, {"price", 1}, {"img_url", 1}, {"user_id", 1}, {"seller_name", 1}, {"latitude", 1}, {"longitude", 1}, {"status_cd", 1}}
	sort := bson.D{{"created_at", -1}}
	opts.SetProjection(projection)
	opts.SetSort(sort)

	// Search User Info
	wg.Add(1)
	go func() {
		profileData.UserInfo = GetSellerInfo(id)
		wg.Done()
	}()

	// Search Seller Items
	wg.Add(1)
	go func() {
		filter := bson.D{{"user_id", id}}
		profileData.Collections = PopulateProductsWithAnImage(filter, opts)
		wg.Done()
	}()
	wg.Wait()

	res.Status = "Success"
	res.Message = "Profile for visitor data retrieved!"
	res.Data = profileData
	json.NewEncoder(c.Writer).Encode(res)
	return
}
