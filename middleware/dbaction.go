package middleware

import (
	"context"
	"log"
	"sync"

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

func PopulateCategories() []interface{} {

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
		appResult.Name = result.Name
		appResult.IconURL = result.IconURL

		results = append(results, appResult)
	}

	if err = cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

func PopulateMainCategories() []interface{} {

	collection := DB.Collection("categories")
	cur, err := collection.Find(context.Background(), bson.D{{"depth", 0}})

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
		appResult.Name = result.Name
		appResult.IconURL = result.IconURL

		results = append(results, appResult)
	}

	if err = cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

func PopulateACategory(id primitive.ObjectID) interface{} {

	var result models.Category

	collection := DB.Collection("categories")
	err := collection.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(&result)

	appResult := struct {
		ID      primitive.ObjectID `json:"_id" bson:"_id"`
		Name    string             `json:"name" bson:"name"`
		IconURL string             `json:"icon_url" bson:"icon_url"`
	}{}

	appResult.ID = result.ID
	appResult.Name = result.Name
	appResult.IconURL = result.IconURL

	if err != nil {
		log.Fatal(err)
	}

	return appResult
}

func FindACategoryFromProductID(id primitive.ObjectID) responses.ProductCategory {
	var query = bson.M{"product_ids": bson.M{"$elemMatch": bson.M{"$eq": id}}}
	var result models.Category

	collection := DB.Collection("categories")
	err := collection.FindOne(context.Background(), query).Decode(&result)
	var appResult responses.ProductCategory

	appResult.ID = result.ID
	appResult.Name = result.Name
	appResult.IconURL = result.IconURL

	if err != nil {
		log.Fatal(err)
	}

	return appResult
}

func GetCategoryIDFromName(categoryName string) primitive.ObjectID {
	columnToSearch := "name"
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

func PopulateRoomUsers(roomID primitive.ObjectID) []interface{} {
	var query = bson.M{"_id": roomID}
	var room models.Room

	collection := DB.Collection("rooms")
	err := collection.FindOne(context.Background(), query).Decode(&room)
	if err != nil {
		log.Fatal(err)
	}

	var results []interface{}
	result := struct {
		ID        primitive.ObjectID `json:"_id" bson:"_id"`
		FirstName string             `json:"first_name" bson:"first_name"`
		LastName  string             `json:"last_name" bson:"last_name"`
		AvatarURL string             `json:"avatar_url" bson:"avatar_url"`
	}{}

	for _, user := range room.UserIDs {
		err := DB.Collection("users").FindOne(context.Background(), bson.M{"_id": user}).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}

	return results
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

	if roomMsg.Image != "" {
		return "An image was sent"
	}
	return roomMsg.Message
}

func PopulateRoomMsgFromRoomID(id primitive.ObjectID, start int64, limit int64) []models.RoomMessage {
	var query = bson.M{"room_id": id}

	limit = limit - start

	if start < 1 {
		start = 1
		limit = 0
	}

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

func PopulateAppointmentsFromUserID(id primitive.ObjectID, userType string) []models.Appointment {

	var collection = DB.Collection("appointments")
	var filter bson.M
	var appointedUserID primitive.ObjectID
	var options = &options.FindOptions{}
	options.SetSort(bson.D{{"created_at", -1}})
	if userType == "seller" {
		filter = bson.M{"seller_id": id}
	} else {
		filter = bson.M{"requester_id": id}
	}

	cur, err := collection.Find(context.Background(), filter, options)
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

		if userType == "seller" {
			appointedUserID = result.RequesterID
		} else {
			appointedUserID = result.SellerID
		}

		result.AppointmentUser = GetUserForNotification(appointedUserID)
		result.ProductDetail = GetAProductWithAnImage(result.ProductID)
		results = append(results, result)
	}

	if results == nil {
		results = make([]models.Appointment, 0)
	}
	return results
}

func PopulateNotificationsFromUserID(idFilter bson.D) []models.Notification {

	var collection = DB.Collection("notifications")
	var opts = &options.FindOptions{}
	opts.SetSort(bson.D{{"created_at", -1}})
	cur, err := collection.Find(context.Background(), idFilter, opts)
	if err != nil {
		log.Fatal(err)
	}

	collection = DB.Collection("appointments")
	var results []models.Notification
	for cur.Next(context.Background()) {
		var result models.Notification
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}

		result.NotificationUser = GetUserForNotification(result.NotifierID)

		if !result.AppointmentID.IsZero() {
			var appointment models.Appointment
			err := collection.FindOne(context.Background(), bson.M{"_id": result.AppointmentID}).Decode(&appointment)
			if err != nil {
				log.Fatal(err)
			}

			result.Appointment = appointment
			result.Product = GetAProductWithAnImage(appointment.ProductID)
		}

		results = append(results, result)
	}

	if results == nil {
		results = make([]models.Notification, 0)
	}
	return results
}

func GetUserForNotification(id primitive.ObjectID) interface{} {
	var collection = DB.Collection("users")

	result := struct {
		ID         primitive.ObjectID `json:"_id" bson:"_id"`
		FirstName  string             `json:"first_name" bson:"first_name"`
		LastName   string             `json:"last_name" bson:"last_name"`
		GoldCoin   int64              `json:"gold_coin"`
		SilverCoin int64              `json:"silver_coin"`
		AvatarURL  string             `json:"avatar_url" bson:"avatar_url"`
	}{}
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func GetSellerInfo(id primitive.ObjectID) interface{} {
	var collection = DB.Collection("users")
	var wg sync.WaitGroup

	result := struct {
		ID         primitive.ObjectID `json:"_id" bson:"_id"`
		FirstName  string             `json:"first_name" bson:"first_name"`
		LastName   string             `json:"last_name" bson:"last_name"`
		CreatedAt  primitive.DateTime `json:"created_at" bson:"created_at"`
		GoldCoin   int64              `json:"gold_coin"`
		SilverCoin int64              `json:"silver_coin"`
		AvatarURL  string             `json:"avatar_url" bson:"avatar_url"`
		ShortBio   string             `json:"short_bio" bson:"short_bio"`
	}{}
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	if result.AvatarURL != "" {
		result.AvatarURL = AvatarURLPrefix + result.ID.Hex() + "/" + result.AvatarURL
	}

	collection = DB.Collection("trust_coins")
	// Search Gold Trust Coins
	wg.Add(1)
	go func() {
		filter := bson.D{bson.E{"receiver_id", id}, bson.E{"type", "gold"}}
		result.GoldCoin, err = collection.CountDocuments(context.Background(), filter)
		wg.Done()
	}()

	// Search Silver Trust Coins
	wg.Add(1)
	go func() {
		filter := bson.D{bson.E{"receiver_id", id}, bson.E{"type", "silver"}}
		result.SilverCoin, err = collection.CountDocuments(context.Background(), filter)
		wg.Done()
	}()
	wg.Wait()

	return result
}

func PopulateCategoriesWithChildren() []interface{} {
	collection := DB.Collection("categories")
	opts := &options.FindOptions{}
	opts.SetSort(bson.D{{"name", 1}})

	cur, err := collection.Find(context.Background(), bson.D{{"depth", 0}}, opts)

	var results []interface{}

	appResult := struct {
		ID       primitive.ObjectID `json:"_id" bson:"_id"`
		Name     string             `json:"name" bson:"name"`
		IconURL  string             `json:"icon_url" bson:"icon_url"`
		Depth    int64              `json:"depth" bson:"depth"`
		Children []interface{}      `json:"children"`
	}{}

	for cur.Next(context.Background()) {
		var result models.Category
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}

		appResult.ID = result.ID
		appResult.Name = result.Name
		appResult.IconURL = result.IconURL
		appResult.Children = getCategoriesChildrenWithRecursion(1, result.ID, 2)

		results = append(results, appResult)
	}

	if err = cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

func getCategoriesChildrenWithRecursion(depth int, parentID primitive.ObjectID, limit int) []interface{} {
	var children []interface{}
	childTemp := struct {
		ID       primitive.ObjectID `json:"_id" bson:"_id"`
		Name     string             `json:"name" bson:"name"`
		IconURL  string             `json:"icon_url" bson:"icon_url"`
		Depth    int64              `json:"depth" bson:"depth"`
		Children []interface{}      `json:"children"`
	}{}

	if depth <= limit {
		var collection = DB.Collection("categories")
		opts := &options.FindOptions{}
		opts.SetSort(bson.D{{"name", 1}})

		childrenCur, err := collection.Find(context.Background(), bson.D{{"parent_id", parentID}, {"depth", depth}}, opts)
		if err != nil {
			log.Fatal(err)
		}
		for childrenCur.Next(context.Background()) {
			var result models.Category
			e := childrenCur.Decode(&result)
			if e != nil {
				log.Fatal(e)
			}

			childTemp.ID = result.ID
			childTemp.Name = result.Name
			childTemp.IconURL = result.IconURL
			childTemp.Depth = result.Depth
			childTemp.Children = getCategoriesChildrenWithRecursion(depth+1, result.ID, limit)
			children = append(children, childTemp)
		}
	}
	if children == nil {
		children = make([]interface{}, 0)
	}
	return children
}
