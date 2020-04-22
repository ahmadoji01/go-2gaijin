package middleware

import (
	"context"
	"fmt"
	"log"

	"gitlab.com/kitalabs/go-2gaijin/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AggregateProductUser(filter bson.D, nPerPage int64, skip int64, sort bson.D) []bson.M {
	matchStage := bson.D{{"$match", filter}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "users"}, {"localField", "user_id"}, {"foreignField", "_id"}, {"as", "user"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$user"}, {"preserveNullAndEmptyArrays", false}}}}
	projectStage := bson.D{{"$project", bson.D{{"_id", 1}, {"name", 1}, {"price", 1}, {"user.first_name", 1}, {"user.email", 1}}}}
	skipStage := bson.D{{"$skip", skip}}
	sortStage := bson.D{{"$sort", sort}}
	limitStage := bson.D{{"$limit", nPerPage}}

	var collection = DB.Collection("products")

	showLoadedCursor, err := collection.Aggregate(context.Background(), mongo.Pipeline{matchStage,
		lookupStage,
		unwindStage,
		projectStage,
		skipStage,
		sortStage,
		limitStage,
	})
	if err != nil {
		panic(err)
	}
	var showsLoaded []bson.M
	if err = showLoadedCursor.All(context.Background(), &showsLoaded); err != nil {
		panic(err)
	}
	return showsLoaded
}

func PopulateProducts(cur *mongo.Cursor, err error) []models.Product {
	var results []models.Product
	for cur.Next(context.Background()) {
		var result models.Product
		fmt.Println(cur)
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

func PopulateBanners(cur *mongo.Cursor, err error) []models.Banner {
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

func PopulateCategories(cur *mongo.Cursor, err error) []models.Category {
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
