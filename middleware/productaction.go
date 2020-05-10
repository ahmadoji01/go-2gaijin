package middleware

import (
	"context"
	"log"
	"strconv"

	"gitlab.com/kitalabs/go-2gaijin/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindAProductImage(id primitive.ObjectID) string {

	result := struct {
		ID    primitive.ObjectID `json:"_id" bson:"_id"`
		Image string             `json:"image" bson:"image"`
	}{}

	coll := DB.Collection("product_images")
	err := coll.FindOne(context.Background(), bson.D{{"product_id", id}}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return ImgURLPrefix + "uploads/product_image/image/" + result.ID.Hex() + "/" + result.Image
}

func PopulateProducts(cur *mongo.Cursor, err error) []models.Product {
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

func PopulateProductsWithAnImage(filter interface{}, options *options.FindOptions) []interface{} {
	var collection = DB.Collection("products")

	cur, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		panic(err)
	}

	result := struct {
		ID         primitive.ObjectID `json:"_id" bson:"_id"`
		Name       string             `json:"name"`
		Price      int                `json:"price"`
		UserID     primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
		SellerName string             `json:"seller_name"`
		ImgURL     string             `json:"img_url"`
		Latitude   string             `json:"latitude,omitempty" bson:"latitude,omitempty"`
		Longitude  string             `json:"longitude,omitempty" bson:"longitude,omitempty"`
		Location   interface{}        `json:"location"`
		StatusEnum int                `json:"status_enum" bson:"status_cd"`
		Status     string             `json:"status" bson:"status"`
	}{}

	var location = struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}{}

	var results []interface{}
	for cur.Next(context.Background()) {
		result.Location = nil
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		result.ImgURL = FindAProductImage(result.ID)
		result.SellerName = FindUserName(result.UserID)
		result.UserID = primitive.NilObjectID

		latitude, e := strconv.ParseFloat(result.Latitude, 64)
		longitude, e := strconv.ParseFloat(result.Longitude, 64)

		location.Latitude = latitude
		location.Longitude = longitude
		result.Location = location

		result.Latitude = ""
		result.Longitude = ""

		result.Status = ProductStatusEnum(result.StatusEnum)

		results = append(results, result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}
