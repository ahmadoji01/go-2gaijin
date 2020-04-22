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
	fmt.Println(payload)
	json.NewEncoder(c.Writer).Encode(payload)
}

func getSearch(query string, category string, nPerPage int64, page int64, sort string, asc int, status int) responses.SearchPage {

	var wg sync.WaitGroup
	var filter bson.D
	var collection = DB.Collection("products")

	filter = searchFilter(query, status)

	var searchData responses.SearchData
	var searchResponse responses.SearchPage
	var count int64
	var err error

	wg.Add(1)
	go func() {
		data := AggregateProductUser(filter, nPerPage, nPerPage*(page-1), bson.D{{sort, asc}})
		searchData.Items = data
		if err != nil {
			searchResponse.Status = "Error"
			searchResponse.Message = err.Error()
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		count, err = collection.CountDocuments(context.Background(), filter)
		if err != nil {
			searchResponse.Status = "Error"
			searchResponse.Message = err.Error()
		}
		wg.Done()
	}()
	wg.Wait()

	if err != nil {
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
		filter = append(filter, bson.E{"name", query})
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
