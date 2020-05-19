package middleware

import (
	"context"
	"log"

	"gitlab.com/kitalabs/go-2gaijin/models"
	"gitlab.com/kitalabs/go-2gaijin/responses"
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

func FindACategoryFromProductID(id primitive.ObjectID, locale string) responses.ProductCategory {
	var query = bson.M{"product_ids": bson.M{"$elemMatch": bson.M{"$eq": id}}}
	var result models.Category

	collection := DB.Collection("categories")
	err := collection.FindOne(context.Background(), query).Decode(&result)

	var appResult responses.ProductCategory

	appResult.ID = result.ID
	appResult.Name = result.Name.Map()[locale].(string)
	appResult.IconURL = result.IconURL

	if err != nil {
		log.Fatal(err)
	}

	return appResult
}

func GetCategoryIDFromName(categoryName string, locale string) primitive.ObjectID {
	columnToSearch := "name." + locale
	query := bson.M{columnToSearch: categoryName}

	collection := DB.Collection("categories")
	result := struct {
		ID primitive.ObjectID `json:"_id" bson:"_id"`
	}{}

	err := collection.FindOne(context.Background(), query).Decode(&result)
	if err != nil {
		return primitive.NilObjectID
	}
	return result.ID
}

func PopulateRoomsFromUserID(id primitive.ObjectID, start int64, limit int64) ([]models.Room, error) {
	var query = bson.M{"user_ids": bson.M{"$elemMatch": bson.M{"$eq": id}}}
	var options = &options.FindOptions{}
	options.SetSort(bson.D{{"last_active", -1}})

	limit = limit - start
	if limit > 16 {
		limit = 16
	}

	if start < 1 {
		start = 1
	}

	options.SetSkip(start - 1)
	options.SetLimit(limit + 1)

	collection := DB.Collection("rooms")
	cur, err := collection.Find(context.Background(), query, options)
	if err != nil {
		log.Fatal(err)
	}

	var results []models.Room
	for cur.Next(context.Background()) {
		var result models.Room
		var anotherUser models.User
		var filter bson.D

		err = cur.Decode(&result)

		if id != result.UserIDs[0] {
			filter = bson.D{{"_id", result.UserIDs[0]}}
		} else {
			filter = bson.D{{"_id", result.UserIDs[1]}}
		}

		collection = DB.Collection("users")
		err = collection.FindOne(context.Background(), filter).Decode(&anotherUser)

		result.Name = anotherUser.FirstName + " " + anotherUser.LastName
		result.IconURL = FindUserAvatar(anotherUser.ID, anotherUser.AvatarURL)
		result.LastMessage = GetLastRoomMsg(result.ID)

		results = append(results, result)
	}

	return results, err
}

func GetLastRoomMsg(id primitive.ObjectID) string {
	var roomMsg models.RoomMessage
	var options = &options.FindOneOptions{}
	options.SetSort(bson.D{{"created_at", -1}})

	collection := DB.Collection("room_messages")
	err := collection.FindOne(context.Background(), bson.D{{"room_id", id}}, options).Decode(&roomMsg)
	if err != nil {
		return ""
	}
	return roomMsg.Message
}

func PopulateRoomMsgFromRoomID(id primitive.ObjectID, start int64, limit int64) []models.RoomMessage {
	var query = bson.M{"room_id": id}

	options := options.Find()
	options.SetSkip(start)
	options.SetLimit(limit)
	options.SetSort(bson.D{{"created_at", 1}})

	collection := DB.Collection("room_messages")
	cur, err := collection.Find(context.Background(), query, options)
	if err != nil {
		log.Fatal(err)
	}

	var results []models.RoomMessage
	for cur.Next(context.Background()) {
		var result models.RoomMessage

		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		results = append(results, result)
	}

	if results == nil {
		results = make([]models.RoomMessage, 0)
	}
	return results
}

func PopulateAppointmentsFromUserID(id primitive.ObjectID) []models.Appointment {

	var collection = DB.Collection("appointments")
	cur, err := collection.Find(context.Background(), bson.M{"requester_id": id})
	if err != nil {
		log.Fatal(err)
	}

	var results []models.Appointment
	for cur.Next(context.Background()) {
		var result models.Appointment

		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		results = append(results, result)
	}

	if results == nil {
		results = make([]models.Appointment, 0)
	}
	return results
}
