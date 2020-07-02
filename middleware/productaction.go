package middleware

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
	"gitlab.com/kitalabs/go-2gaijin/config"
	"gitlab.com/kitalabs/go-2gaijin/models"
	"gitlab.com/kitalabs/go-2gaijin/responses"
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
	if !strings.HasPrefix(result.Image, "https://") {
		return ImgURLPrefix + "uploads/product_image/image/" + result.ID.Hex() + "/" + result.Image
	}
	return result.Image
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
		ID           primitive.ObjectID `json:"_id" bson:"_id"`
		Name         string             `json:"name"`
		Price        int                `json:"price"`
		UserID       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
		SellerName   string             `json:"seller_name"`
		ImgURL       string             `json:"img_url"`
		Latitude     float64            `json:"latitude,omitempty" bson:"latitude,omitempty"`
		Longitude    float64            `json:"longitude,omitempty" bson:"longitude,omitempty"`
		Location     interface{}        `json:"location"`
		IsLiked      bool               `json:"is_liked"`
		StatusEnum   int                `json:"status_cd" bson:"status_cd"`
		Availability string             `json:"availability" bson:"availability"`
	}{}

	var location = struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}{}

	var results []interface{}
	for cur.Next(context.Background()) {
		result.Location = nil
		var prod models.Product
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		result.ImgURL = FindAProductImage(result.ID)
		result.SellerName = FindUserName(result.UserID)

		err := collection.FindOne(context.Background(), bson.M{"$and": bson.D{{"_id", result.ID}, {"liked_by", result.UserID}}}).Decode(&prod)
		if err == nil {
			result.IsLiked = true
		} else {
			result.IsLiked = false
		}

		result.UserID = primitive.NilObjectID
		location.Latitude = result.Latitude
		location.Longitude = result.Longitude
		result.Location = location

		result.Latitude = 0
		result.Longitude = 0

		result.Availability = ProductStatusEnum(result.StatusEnum)

		results = append(results, result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())

	if results == nil {
		results = make([]interface{}, 0)
	}

	return results
}

func GetAProductWithAnImage(id primitive.ObjectID) interface{} {
	var collection = DB.Collection("products")
	var product models.Product

	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&product)
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
		Location   interface{}        `json:"location"`
		IsLiked    bool               `json:"is_liked"`
		StatusEnum int                `json:"status_enum" bson:"status_cd"`
		Status     string             `json:"status" bson:"status"`
	}{}

	var location = struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}{}

	var prod models.Product
	err = collection.FindOne(context.Background(), bson.M{"$and": bson.D{{"_id", result.ID}, {"liked_by", result.UserID}}}).Decode(&prod)
	if err == nil {
		result.IsLiked = true
	} else {
		result.IsLiked = false
	}

	result.ID = product.ID
	result.Name = product.Name
	result.Price = product.Price
	result.UserID = product.User
	result.SellerName = FindUserName(product.User)
	result.ImgURL = FindAProductImage(product.ID)
	location.Latitude = product.Latitude
	location.Longitude = product.Longitude
	result.Location = location
	result.Status = ProductStatusEnum(product.StatusEnum)
	result.StatusEnum = product.StatusEnum

	return result
}

func DeleteProduct(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var product models.Product
	var res responses.GenericResponse

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &product)
	if err != nil {
		res.Status = "Error"
		res.Message = "Something wrong happened. Try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)

	if isLoggedIn {
		var collection = DB.Collection("products")
		err := collection.FindOne(context.Background(), bson.M{"_id": product.ID}).Decode(&product)
		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		if product.User != userData.ID {
			res.Status = "Error"
			res.Message = "You are not authorized to delete this product"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		collection = DB.Collection("products")
		_, err = collection.DeleteOne(context.Background(), bson.D{{"_id", product.ID}})
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		res.Status = "Success"
		res.Message = "You have deleted this product"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func EditProduct(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")
	var res responses.GenericResponse

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)

	if isLoggedIn {
		isSubscribed := IsUserSubscribed(userData.ID)
		isSubscribed = true

		if isSubscribed {
			var productInsert models.ProductInsert
			var productData models.Product

			body, _ := ioutil.ReadAll(c.Request.Body)
			err := json.Unmarshal(body, &productInsert)
			if err != nil {
				res.Status = "Error"
				res.Message = err.Error()
				json.NewEncoder(c.Writer).Encode(res)
				return
			}

			err = DB.Collection("products").FindOne(context.Background(), bson.M{"_id": productInsert.Product.ID}).Decode(&productData)
			if err != nil {
				res.Status = "Error"
				res.Message = "Product not found!"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}

			if productData.User != userData.ID {
				res.Status = "Error"
				res.Message = "You are not authorized to edit this product"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}

			var geoLoc models.GeoJson
			geoLoc.Type = "Point"
			geoLoc.Coordinates = append(geoLoc.Coordinates, productInsert.Product.Longitude)
			geoLoc.Coordinates = append(geoLoc.Coordinates, productInsert.Product.Latitude)

			var collection = DB.Collection("products")
			update := bson.M{"$set": bson.M{
				"name":         productInsert.Product.Name,
				"price":        productInsert.Product.Price,
				"description":  productInsert.Product.Description,
				"category_ids": productInsert.Product.Category,
				"latitude":     productInsert.Product.Latitude,
				"longitude":    productInsert.Product.Longitude,
				"geoloc":       geoLoc,
			}}
			_, err = collection.UpdateOne(context.Background(), bson.M{"_id": productInsert.Product.ID}, update)
			if err != nil {
				res.Status = "Error"
				res.Message = "Something went wrong"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}

			collection = DB.Collection("product_details")
			update = bson.M{"$set": bson.M{
				"brand":       productInsert.ProductDetail.Brand,
				"condition":   productInsert.ProductDetail.Condition,
				"years_owned": productInsert.ProductDetail.YearsOwned,
				"model_name":  productInsert.ProductDetail.ModelName,
			}}
			_, err = collection.UpdateOne(context.Background(), bson.M{"product_id": productInsert.Product.ID}, update)
			if err != nil {
				res.Status = "Error"
				res.Message = "Something went wrong"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}

			res.Status = "Success"
			res.Message = "Product Successfully Edited"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}
		res.Status = "Error"
		res.Message = "Unsubscribed"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func GetProductInfoForEdit(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")
	var res responses.GenericResponse

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)

	var product models.Product
	var productDetail models.ProductDetail

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &product)
	if err != nil {
		res.Status = "Error"
		res.Message = "Something wrong happened. Try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	if isLoggedIn {
		var collection = DB.Collection("products")
		err := collection.FindOne(context.Background(), bson.M{"_id": product.ID}).Decode(&product)
		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		if product.User != userData.ID {
			res.Status = "Error"
			res.Message = "You are not authorized to get the product information"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		collection = DB.Collection("product_details")
		err = collection.FindOne(context.Background(), bson.M{"product_id": product.ID}).Decode(&productDetail)
		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		var productEditInfo responses.ProductEditInfo
		productEditInfo.Product = product
		productEditInfo.ProductDetail = productDetail

		res.Status = "Success"
		res.Message = "Product Information for Editing Successfully Retrieved"
		res.Data = productEditInfo
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func LikeProduct(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var product models.Product
	var res responses.GenericResponse

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &product)
	if err != nil {
		res.Status = "Error"
		res.Message = "Something wrong happened. Try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		var prod models.Product
		var collection = DB.Collection("products")

		err := collection.FindOne(context.Background(), bson.M{"$and": bson.D{{"_id", product.ID}, {"liked_by", userData.ID}}}).Decode(&prod)
		if err == nil {
			res.Status = "Error"
			res.Message = "Product is already liked"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		update := bson.M{"$push": bson.M{"liked_by.$.arr": userData.ID}}
		_, err = collection.UpdateOne(context.Background(), bson.M{"_id": product.ID}, update)
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}
		res.Status = "Success"
		res.Message = "Product has successfully been liked"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func MarkAsSold(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var product models.Product
	var res responses.GenericResponse

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &product)
	if err != nil {
		res.Status = "Error"
		res.Message = "Something wrong happened. Try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)

	if isLoggedIn {
		var collection = DB.Collection("products")
		err := collection.FindOne(context.Background(), bson.M{"_id": product.ID}).Decode(&product)
		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		if product.User != userData.ID {
			res.Status = "Error"
			res.Message = "You are not authorized to mark this product as sold"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		collection = DB.Collection("products")

		var update bson.M
		if product.Availability == "available" {
			update = bson.M{"$set": bson.M{"availability": "sold", "status_cd": 2}}
			res.Message = "You have marked this product as sold"
		} else {
			update = bson.M{"$set": bson.M{"availability": "available", "status_cd": 1}}
			res.Message = "You have marked this product as available"
		}

		_, err = collection.UpdateOne(context.Background(), bson.D{{"_id", product.ID}}, update)
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		res.Status = "Success"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func EditPricing(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var product models.Product
	var res responses.GenericResponse

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &product)
	if err != nil {
		res.Status = "Error"
		res.Message = "Something wrong happened. Try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)

	if isLoggedIn {
		price := product.Price
		var collection = DB.Collection("products")
		err := collection.FindOne(context.Background(), bson.M{"_id": product.ID}).Decode(&product)
		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		if product.User != userData.ID {
			res.Status = "Error"
			res.Message = "You are not authorized to edit the price for this product"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		collection = DB.Collection("products")
		update := bson.M{"$set": bson.M{"price": price}}

		_, err = collection.UpdateOne(context.Background(), bson.D{{"_id", product.ID}}, update)
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		res.Status = "Success"
		res.Message = "Price has successfully been changed"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func GetAllCategories(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	catData := struct {
		Categories interface{} `json:"categories"`
	}{}

	var res responses.GenericResponse

	catData.Categories = PopulateCategoriesWithChildren()

	res.Status = "Success"
	res.Message = "Categories Retrieved"
	res.Data = catData
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func PostNewProduct(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")
	var res responses.GenericResponse

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)

	if isLoggedIn {
		isSubscribed := IsUserSubscribed(userData.ID)
		isSubscribed = true

		if isSubscribed {
			var productInsert models.ProductInsert

			body, _ := ioutil.ReadAll(c.Request.Body)
			err := json.Unmarshal(body, &productInsert)
			if err != nil {
				res.Status = "Error"
				res.Message = err.Error()
				json.NewEncoder(c.Writer).Encode(res)
				return
			}

			var geoLoc models.GeoJson
			geoLoc.Type = "Point"
			geoLoc.Coordinates = append(geoLoc.Coordinates, productInsert.Product.Longitude)
			geoLoc.Coordinates = append(geoLoc.Coordinates, productInsert.Product.Latitude)

			productInsert.ProductDetail.ID = primitive.NewObjectIDFromTimestamp(time.Now())
			productInsert.Product.ID = primitive.NewObjectIDFromTimestamp(time.Now())
			productInsert.Product.StatusEnum = 1
			productInsert.Product.Availability = "available"

			productInsert.Product.User = userData.ID
			productInsert.Product.DateCreated = primitive.NewDateTimeFromTime(time.Now())
			productInsert.Product.DateUpdated = primitive.NewDateTimeFromTime(time.Now())
			productInsert.Product.ProductDetails = productInsert.ProductDetail.ID
			productInsert.Product.GeoLoc = geoLoc

			var collection = DB.Collection("products")
			productData, err := collection.InsertOne(context.Background(), productInsert.Product)
			if err != nil {
				res.Status = "Error"
				res.Message = "Something went wrong"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}
			uploadProductImages(productData.InsertedID.(primitive.ObjectID), productInsert.ProductImages)

			collection = DB.Collection("product_details")
			productInsert.ProductDetail.ProductID = productData.InsertedID.(primitive.ObjectID)
			_, err = collection.InsertOne(context.Background(), productInsert.ProductDetail)
			if err != nil {
				res.Status = "Error"
				res.Message = "Something went wrong"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}

			dir, err := ioutil.ReadDir("tmp")
			for _, d := range dir {
				os.RemoveAll(path.Join([]string{"tmp", d.Name()}...))
			}

			res.Status = "Success"
			res.Message = "Product Successfully Inserted"
			res.Data = productData
			json.NewEncoder(c.Writer).Encode(res)
			return
		}
		res.Status = "Error"
		res.Message = "Unsubscribed"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func uploadProductImages(productID primitive.ObjectID, productImages []models.ProductImage) {
	var collection = DB.Collection("product_images")
	var wg sync.WaitGroup

	for i := 0; i < len(productImages); i++ {
		imgPath := uuid.NewV4().String()
		imgPath = imgPath + "/"

		imgName := uuid.NewV4().String()
		imgName = imgName + ".jpg"

		imgData := productImages[i].ImgData
		thumbData := productImages[i].ThumbData

		productImages[i].Order = i + 1
		productImages[i].Product = productID
		productImages[i].ImgURL = GCSProductImgPrefix + imgPath + imgName
		productImages[i].ThumbURL = GCSProductImgPrefix + imgPath + "thumb_" + imgName
		productImages[i].ImgData = ""
		productImages[i].ThumbData = ""

		_, err := collection.InsertOne(context.Background(), productImages[i])
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		go func() {
			DecodeBase64ToImage(imgData, imgName)
			UploadToGCS(ProductImagePrefix+imgPath, imgName)
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			DecodeBase64ToImage(thumbData, "thumb_"+imgName)
			UploadToGCS(ProductImagePrefix+imgPath, "thumb_"+imgName)
			wg.Done()
		}()
		wg.Wait()
	}
}
