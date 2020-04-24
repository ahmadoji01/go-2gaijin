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
		ID     primitive.ObjectID `json:"_id" bson:"_id"`
		Name   string             `json:"name"`
		Price  int                `json:"price"`
		ImgURL string             `json:"img_url"`
	}{}

	var results []interface{}
	for cur.Next(context.Background()) {
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		result.ImgURL = ImgURLPrefix + "uploads/product_image/image/" + result.ID.Hex() + "/" + FindAProductImage(result.ID)
		results = append(results, result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

func FindAProductImage(id primitive.ObjectID) string {
	result := struct {
		Image string
	}{}

	coll := DB.Collection("product_images")
	err := coll.FindOne(context.Background(), bson.D{{"product_id", id}}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result.Image
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
