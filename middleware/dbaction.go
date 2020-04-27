package middleware

import (
	"context"
	"log"

	"gitlab.com/kitalabs/go-2gaijin/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndex(weights bson.M, keys bson.M, coll *mongo.Collection) {
	opt := options.Index()
	opt.SetWeights(weights)

	index := mongo.IndexModel{Keys: keys, Options: opt}
	if _, err := coll.Indexes().CreateOne(context.Background(), index); err != nil {
		log.Println("Could not create text index:", err)
	}
}

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

func SearchProducts(filter bson.D, options *options.FindOptions) []bson.M {
	var collection = DB.Collection("products")

	cur, err := collection.Find(context.Background(), filter, options)
	if err != nil {
		panic(err)
	}

	var showsLoaded []bson.M
	if err = cur.All(context.Background(), &showsLoaded); err != nil {
		panic(err)
	}
	return showsLoaded
}

func PopulateModels(cur *mongo.Cursor, err error) []bson.M {
	var results []bson.M
	for cur.Next(context.Background()) {
		var result bson.M
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
		UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
		SellerName string             `json:"seller_name"`
		ImgURL     string             `json:"img_url"`
		Location   []float64          `json:"location" bson:"location"`
		StatusEnum int                `json:"status_enum" bson:"status_cd"`
		Status     string             `json:"status" bson:"status"`
	}{}

	var results []interface{}
	for cur.Next(context.Background()) {
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		result.ImgURL = FindAProductImage(result.ID)
		result.SellerName = FindUserName(result.UserID)
		result.UserID = primitive.NilObjectID

		if result.StatusEnum == 1 {
			result.Status = "available"
		} else if result.StatusEnum == 2 {
			result.Status = "sold"
		} else {
			result.Status = "unavailable"
		}

		results = append(results, result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

func FindUserName(id primitive.ObjectID) string {
	result := struct {
		FirstName string `json:"first_name" bson:"first_name"`
		LastName  string `json:"last_name" bson:"last_name"`
	}{}

	coll := DB.Collection("users")
	err := coll.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result.FirstName
}

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

func FindProductImages(id string) []bson.M {
	var results []bson.M

	productID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}

	coll := DB.Collection("product_images")
	cur, err := coll.Find(context.Background(), bson.D{{"product_id", productID}})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		results = append(results, result)
	}
	return results
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

func PopulateCategories(locale string) []interface{} {

	collection := DB.Collection("categories")
	cur, err := collection.Find(context.Background(), bson.D{{}})

	var results []interface{}

	appResult := struct {
		ID      primitive.ObjectID `json:"_id" bson:"_id"`
		Name    string             `json:"name" bson:"name"`
		IconURL string             `json:"icon_url" bson:"icon_url"`
	}{}

	for cur.Next(context.Background()) {
		var result models.Category
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}

		appResult.ID = result.ID
		appResult.Name = result.Name.Map()[locale].(string)
		appResult.IconURL = result.IconURL

		results = append(results, appResult)
	}

	if err = cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

func PopulateACategory(id primitive.ObjectID, locale string) interface{} {

	var result models.Category

	collection := DB.Collection("categories")
	err := collection.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(&result)

	appResult := struct {
		ID      primitive.ObjectID `json:"_id" bson:"_id"`
		Name    string             `json:"name" bson:"name"`
		IconURL string             `json:"icon_url" bson:"icon_url"`
	}{}

	appResult.ID = result.ID
	appResult.Name = result.Name.Map()[locale].(string)
	appResult.IconURL = result.IconURL

	if err != nil {
		log.Fatal(err)
	}

	return appResult
}

func FindACategoryFromProductID(id primitive.ObjectID, locale string) interface{} {
	var query = bson.M{"product_ids": bson.M{"$elemMatch": bson.M{"$eq": id}}}
	var result models.Category

	collection := DB.Collection("categories")
	err := collection.FindOne(context.Background(), query).Decode(&result)

	appResult := struct {
		ID      primitive.ObjectID `json:"_id" bson:"_id"`
		Name    string             `json:"name" bson:"name"`
		IconURL string             `json:"icon_url" bson:"icon_url"`
	}{}

	appResult.ID = result.ID
	appResult.Name = result.Name.Map()[locale].(string)
	appResult.IconURL = result.IconURL

	if err != nil {
		log.Fatal(err)
	}

	return appResult
}
