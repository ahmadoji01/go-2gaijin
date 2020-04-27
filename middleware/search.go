package middleware

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/kitalabs/go-2gaijin/models"
	"gitlab.com/kitalabs/go-2gaijin/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetSearch(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

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

	var status = -1
	if urlQuery.Get("status") == "" {
		status = -1
	} else {
		status, err = strconv.Atoi(urlQuery.Get("status"))
	}

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

	var res models.ResponseResult
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	payload := getSearch(query, category, start, limit, priceMin, priceMax, sort, asc, status)
	json.NewEncoder(c.Writer).Encode(payload)
}

func getSearch(query string, category string, start int64, limit int64, priceMin int64, priceMax int64, sort string, asc int, status int) interface{} {

	//filter := bson.D{bson.E{"$text", bson.M{"$search": query}}}
	filter := searchFilter(query, status, priceMin, priceMax)
	findOptions := options.Find()
	findOptions.SetLimit(limit - start)
	findOptions.SetSkip(start)
	findOptions.SetProjection(bson.M{
		"_id":         1,
		"name":        1,
		"price":       1,
		"description": 1,
		"user_id":     1,
		"location":    1,
		"status_cd":   1,
		"relevance":   bson.M{"$meta": "textScore"},
	})
	findOptions.SetSort(bson.M{"relevance": bson.M{"$meta": "textScore"}})

	return PopulateProductsWithAnImage(filter, findOptions)
}

func searchFilter(query string, status int, priceMin int64, priceMax int64) bson.D {

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
