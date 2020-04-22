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

func GetSearch(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	urlQuery := c.Request.URL.Query()

	query := urlQuery.Get("q")
	sort := urlQuery.Get("sortby")
	if sort == "" {
		sort = "created_at"
	}

	var err error
	asc := -1
	if urlQuery.Get("asc") != "" {
		asc, err = strconv.Atoi(urlQuery.Get("asc"))
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
		if page <= 0 {
			page = 1
		}
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

func getSearch(query string, category string, nPerPage int64, page int64, sort string, asc int, status int) responses.SearchPage {

	var wg sync.WaitGroup
	var filter bson.D
	var collection = DB.Collection("products")

	filter = searchFilter(query, status)
	fmt.Println(filter)

	var options = &options.FindOptions{}
	options.SetLimit(nPerPage)
	options.SetSort(bson.D{{sort, asc}})
	options.SetSkip(nPerPage * (page - 1))

	var searchData responses.SearchData
	var searchResponse responses.SearchPage
	var count int64
	var err error

	wg.Add(1)
	go func() {
		searchData.Items = PopulateProducts(collection.Find(context.Background(), filter, options))
		wg.Done()
	}()

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

	searchData.Pagination = getPagination(count, nPerPage, page)

	searchResponse.Data = searchData
	searchResponse.Status = "Success"
	searchResponse.Message = "Search Page Data Loaded"

	return searchResponse
}

func searchFilter(query string, status int) bson.D {
	var filter = bson.D{}

	if query != "" {
		filter = append(filter, bson.E{"_keywords", bson.M{"$elemMatch": bson.M{"$regex": "/*." + query + ".*/"}}})
	}
	if status != -1 {
		filter = append(filter, bson.E{"status_cd", status})
	}

	return filter

}

func getPagination(totalItems int64, nPerPage int64, currentPage int64) responses.Pagination {
	var pagination responses.Pagination
	pagination.CurrentPage = currentPage
	pagination.ItemsPerPage = nPerPage
	pagination.TotalItems = totalItems

	if totalItems-(nPerPage*(currentPage+1)) >= 1 {
		pagination.NextPage = currentPage + 1
	} else {
		pagination.NextPage = 0
	}

	if currentPage-1 >= 1 {
		pagination.PreviousPage = currentPage - 1
	} else {
		pagination.PreviousPage = 0
	}

	pagination.TotalPages = totalItems / nPerPage

	return pagination
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
