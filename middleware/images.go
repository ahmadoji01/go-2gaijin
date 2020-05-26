package middleware

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Image URL Prefix
var ImgURLPrefix string = "https://storage.googleapis.com/rails-2gaijin-storage/"
var AvatarURLPrefix string = "https://storage.googleapis.com/rails-2gaijin-storage/uploads/user/avatar/"

func FindProductImages(productID primitive.ObjectID) []interface{} {
	var results []interface{}

	coll := DB.Collection("product_images")
	cur, err := coll.Find(context.Background(), bson.D{{"product_id", productID}})
	if err != nil {
		log.Fatal(err)
	}

	var result = struct {
		ID    primitive.ObjectID `json:"_id" bson:"_id"`
		Image string             `json:"img_url" bson:"image"`
	}{}

	for cur.Next(context.Background()) {
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		result.Image = ImgURLPrefix + "uploads/product_image/image/" + result.ID.Hex() + "/" + result.Image

		results = append(results, result)
	}
	return results
}

func FindUserAvatar(userID primitive.ObjectID, avatarName string) string {
	if avatarName == "" {
		return ""
	} else {
		var avatarURL = ImgURLPrefix + "uploads/user/avatar/" + userID.Hex() + "/" + avatarName
		return avatarURL
	}
}
