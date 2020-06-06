package middleware

import (
	"context"
	"encoding/base64"
	"io"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/option"
)

//Image URL Prefix
var ImgURLPrefix string = "https://storage.googleapis.com/rails-2gaijin-storage/"
var AvatarURLPrefix string = "https://storage.googleapis.com/rails-2gaijin-storage/uploads/user/avatar/"
var ProductImagePrefix = "uploads/go_product_image/"
var GCSProductImgPrefix = "https://storage.googleapis.com/rails-2gaijin-storage/uploads/go_product_image/"

// Authenticate to Google Cloud Storage and return handler
func UploadToGCS(filePath string, fileName string) {

	credentialFilePath := "./keys/rails-2gaijin-790de45ba7c6.json"

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialFilePath))
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open("tmp/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	bucketName := "rails-2gaijin-storage"
	objectPath := filePath + fileName
	obj := client.Bucket(bucketName).Object(objectPath)
	wc := obj.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		log.Fatal(err)
	}
	if err := wc.Close(); err != nil {
		log.Fatal(err)
	}
	log.Println("done")
}

func DecodeBase64ToImage(str string, filename string) *os.File {
	data, _ := base64.StdEncoding.DecodeString(str) //[]byte

	file, _ := os.Create("tmp/" + filename)
	defer file.Close()

	file.Write(data)
	return file
}

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

		if !strings.HasPrefix(result.Image, "https://") {
			result.Image = ImgURLPrefix + "uploads/product_image/image/" + result.ID.Hex() + "/" + result.Image
		}
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
