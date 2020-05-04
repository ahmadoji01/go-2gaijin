package middleware

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/kitalabs/go-2gaijin/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetSearch(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	urlQuery := c.Request.URL.Query()

	query := urlQuery.Get("q")
	sort := urlQuery.Get("sortby")
	if sort == "" {
		sort = "relevance"
	}

	var err error
	asc := -1
	if urlQuery.Get("asc") != "" {
		asc, err = strconv.Atoi(urlQuery.Get("asc"))
	}

	var status = urlQuery.Get("status")

	category := urlQuery.Get("category")
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

	var priceMin int64
	var priceMax int64

	if urlQuery.Get("pricemin") == "" {
		priceMin = -1
	} else {
		priceMin, err = strconv.ParseInt(urlQuery.Get("pricemin"), 10, 64)
	}
	if urlQuery.Get("pricemax") == "" {
		priceMax = -1
	} else {
		priceMax, err = strconv.ParseInt(urlQuery.Get("pricemax"), 10, 64)
	}

	var res responses.ResponseMessage
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	payload := getSearch(query, category, start, limit, priceMin, priceMax, sort, asc, status)

	var searchPage responses.SearchPage
	searchPage.Status = "Success"
	searchPage.Message = "Products Successfully Searched"
	searchPage.Data = payload

	json.NewEncoder(c.Writer).Encode(searchPage)
}

func getSearch(query string, category string, start int64, limit int64, priceMin int64, priceMax int64, sort string, asc int, status string) interface{} {

	filter := searchFilter(query, status, priceMin, priceMax, category)
	findOptions := searchOptions(start, limit, sort)
	findOptions.SetProjection(bson.M{
		"_id":         1,
		"name":        1,
		"price":       1,
		"description": 1,
		"user_id":     1,
		"latitude":    1,
		"longitude":   1,
		"location":    1,
		"status_cd":   1,
		"relevance":   bson.M{"$meta": "textScore"},
	})

	return PopulateProductsWithAnImage(filter, findOptions)
}

func searchFilter(query string, status string, priceMin int64, priceMax int64, category string) bson.D {

	var filter bson.D
	if priceMax != -1 && priceMin != -1 {
		filter = append(filter, bson.E{"price", bson.M{"$lte": priceMax, "$gte": priceMin}})
	} else if priceMin != -1 && priceMax == -1 {
		filter = append(filter, bson.E{"price", bson.M{"$gte": priceMin}})
	} else if priceMax != -1 && priceMin == -1 {
		filter = append(filter, bson.E{"price", bson.M{"$lte": priceMax}})
	}

	if query != "" {
		filter = append(filter, bson.E{"$text", bson.M{"$search": query}})
	}

	if status == "sold" {
		filter = append(filter, bson.E{"status_cd", 2})
	} else if status == "available" {
		filter = append(filter, bson.E{"status_cd", 1})
	}

	if category != "" {
		cat := GetCategoryIDFromName(category, "en")
		if cat != primitive.NilObjectID {
			filter = append(filter, bson.E{"category_ids", cat})
		}
	}

	return filter
}

func searchOptions(start int64, limit int64, sort string) *options.FindOptions {
	findOptions := options.Find()

	limit = limit - start
	if limit > 16 {
		limit = 16
	}

	findOptions.SetSkip(start)
	findOptions.SetLimit(limit)

	if sort == "relevance" {
		findOptions.SetSort(bson.M{"relevance": bson.M{"$meta": "textScore"}})
	} else if sort == "newest" {
		findOptions.SetSort(bson.D{{"created_at", 1}})
	} else if sort == "oldest" {
		findOptions.SetSort(bson.D{{"created_at", -1}})
	} else if sort == "highestprice" {
		findOptions.SetSort(bson.D{{"price", -1}})
	} else if sort == "lowestprice" {
		findOptions.SetSort(bson.D{{"price", 1}})
	}

	return findOptions
}

func getPagination(totalItems int64, nPerPage int64, currentPage int64) responses.Pagination {
	var pagination responses.Pagination
	pagination.CurrentPage = currentPage
	pagination.ItemsPerPage = nPerPage
	pagination.TotalItems = totalItems

	if totalItems-(nPerPage*currentPage) >= 1 {
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
