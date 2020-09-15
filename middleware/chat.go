package middleware

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
	"gitlab.com/kitalabs/go-2gaijin/channels"
	"gitlab.com/kitalabs/go-2gaijin/config"
	"gitlab.com/kitalabs/go-2gaijin/models"
	"gitlab.com/kitalabs/go-2gaijin/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetChatRoomMsg(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var urlQuery = c.Request.URL.Query()
	collection := DB.Collection("room_messages")

	var msgsData []models.RoomMessage
	var roomData responses.ChatRoomData
	var res responses.GenericResponse

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		var start, limit int64

		id, err := primitive.ObjectIDFromHex(urlQuery.Get("room"))
		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		totalMessages, _ := collection.CountDocuments(context.Background(), bson.D{{"room_id", id}})

		if urlQuery.Get("limit") == "" {
			limit = totalMessages
		} else {
			limit, err = strconv.ParseInt(urlQuery.Get("limit"), 10, 64)
		}

		if urlQuery.Get("start") == "" {
			start = limit - 24
			if start <= 0 {
				start = 0
			}
		} else {
			start, err = strconv.ParseInt(urlQuery.Get("start"), 10, 64)
		}
		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		msgsData = PopulateRoomMsgFromRoomID(id, start, limit)
		roomData.Messages = msgsData
		roomData.TotalMessages = totalMessages
		setRoomMsgReader(id, userData.ID)

		res.Status = "Success"
		res.Message = "Chat Messages Successfully Retrieved"
		res.Data = roomData
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func GetChatRoomUser(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var urlQuery = c.Request.URL.Query()

	var roomUsersData responses.ChatRoomUsersData
	var res responses.GenericResponse

	tokenString := c.Request.Header.Get("Authorization")
	_, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		id, err := primitive.ObjectIDFromHex(urlQuery.Get("room"))
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		roomUsersData.Users = PopulateRoomUsers(id)

		res.Status = "Success"
		res.Message = "Room's Users Data Successfully Retrieved"
		res.Data = roomUsersData
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func GetChatLobby(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Content-Type", "application/json")

	var res responses.GenericResponse
	var err error

	urlQuery := c.Request.URL.Query()
	var start int64
	var limit int64
	if urlQuery.Get("start") == "" {
		start = 0
	} else {
		start, err = strconv.ParseInt(urlQuery.Get("start"), 10, 64)
	}

	if urlQuery.Get("limit") == "" {
		limit = 24
	} else {
		limit, err = strconv.ParseInt(urlQuery.Get("limit"), 10, 64)
		if limit <= 0 {
			limit = 24
		}
	}

	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	var roomsData []models.Room
	var lobbyData responses.ChatLobbyData

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		_, err = DB.Collection("users").UpdateOne(context.Background(), bson.M{"_id": userData.ID}, bson.M{"$set": bson.M{"message_read": true}})
		if err != nil {
			log.Fatal(err)
		}

		roomsData, _ = PopulateRoomsFromUserID(userData.ID, start, limit)
		res.Message = "Chat Lobby Retrieved!"
		res.Status = "Success"
		if roomsData == nil {
			roomsData = make([]models.Room, 0)
		}

		lobbyData.ChatLobby = roomsData
		res.Data = lobbyData
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func InsertMessage(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var roomMsg models.RoomMessage
	var res responses.GenericResponse

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &roomMsg)
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	roomMsg.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	roomMsg.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	var collection = DB.Collection("room_messages")
	tokenString := c.Request.Header.Get("Authorization")
	_, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		_, err := collection.InsertOne(context.TODO(), roomMsg)

		collection = DB.Collection("rooms")
		_, err = collection.UpdateOne(context.Background(), bson.M{"_id": roomMsg.RoomID}, bson.D{{"$set", bson.D{{"last_active", roomMsg.CreatedAt}}}})

		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		notifyUnreadMessage(roomMsg.UserID)

		var roomMsgRes responses.InsertRoomMsg
		roomMsgRes.RoomMsg = roomMsg

		res.Status = "Success"
		res.Message = "Message successfully saved"
		res.Data = roomMsgRes
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func InsertImageMessage(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var roomMsg models.RoomMessage
	var res responses.GenericResponse

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &roomMsg)
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	roomMsg.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	roomMsg.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	var collection = DB.Collection("room_messages")
	tokenString := c.Request.Header.Get("Authorization")
	_, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		if roomMsg.ImgData == "" {
			res.Status = "Error"
			res.Message = "No image was attached"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		imgName := uuid.NewV4().String()
		imgName = imgName + ".jpg"

		imgPath := GCSChatImgPrefix + roomMsg.RoomID.Hex() + "/"

		DecodeBase64ToImage(roomMsg.ImgData, imgName)
		UploadToGCS(ChatImagePrefix+roomMsg.RoomID.Hex()+"/", imgName)

		roomMsg.Message = ""
		roomMsg.ImgData = ""
		roomMsg.Image = imgPath + imgName

		_, err := collection.InsertOne(context.TODO(), roomMsg)

		collection = DB.Collection("rooms")
		_, err = collection.UpdateOne(context.Background(), bson.M{"_id": roomMsg.RoomID}, bson.D{{"$set", bson.D{{"last_active", roomMsg.CreatedAt}}}})

		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		notifyUnreadMessage(roomMsg.UserID)

		var roomMsgImg responses.InsertRoomMsg
		roomMsgImg.RoomMsg = roomMsg

		res.Status = "Success"
		res.Message = "Image Message Successfully Saved"
		res.Data = roomMsgImg
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func notifyUnreadMessage(userID primitive.ObjectID) {
	var collection = DB.Collection("users")
	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": userID}, bson.M{"$set": bson.M{"message_read": false}})
	if err != nil {
		log.Fatal(err)
	}

	msg := `{"message_read": false}`
	msgByte := []byte(msg)

	var m channels.Message

	m.Data = msgByte
	m.Room = userID.Hex()
	channels.H.Broadcast <- m
}

func setIsRead(roomID primitive.ObjectID) {
	var collection = DB.Collection("rooms")
	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": roomID}, bson.M{"$set": bson.M{"is_read": true}})
	if err != nil {
		log.Fatal(err)
	}
}

func setRoomMsgReader(roomID primitive.ObjectID, userID primitive.ObjectID) {
	var collection = DB.Collection("room_messages")

	roomMsg := GetLastRoomMsg(roomID)

	if !IDExistsInSlice(roomMsg.ReaderIDs, userID) {
		readerIDs := append(roomMsg.ReaderIDs, userID)
		_, err := collection.UpdateOne(context.Background(), bson.M{"_id": roomMsg.ID}, bson.M{"$set": bson.M{"reader_ids": readerIDs}})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ChatUser(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var room models.Room
	var res responses.GenericResponse
	collection := DB.Collection("rooms")

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		receiverID, _ := primitive.ObjectIDFromHex(c.Request.URL.Query().Get("receiverid"))

		query := bson.D{{"user_ids", bson.D{{"$all", bson.A{userData.ID, receiverID}}}}}
		roomNotFound := collection.FindOne(context.Background(), query).Decode(&room)
		if roomNotFound != nil {
			room.LastActive = primitive.NewDateTimeFromTime(time.Now())
			room.UserIDs = []primitive.ObjectID{userData.ID, receiverID}
			result, err := collection.InsertOne(context.TODO(), room)
			if err != nil {
				res.Status = "Error"
				res.Message = "Something wrong happened. Try again"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}

			room.ID = result.InsertedID.(primitive.ObjectID)
		}
		var roomData responses.RoomData
		roomData.Room = room

		res.Status = "Success"
		res.Message = "Chat room has been obtained"
		res.Data = roomData
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func AddMsgReader(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var roomMsg models.RoomMessage
	var res responses.GenericResponse
	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &roomMsg)
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	collection := DB.Collection("room_messages")

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		err = collection.FindOne(context.Background(), bson.M{"_id": roomMsg.ID}).Decode(&roomMsg)
		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		if !IDExistsInSlice(roomMsg.ReaderIDs, userData.ID) {
			readerIDs := append(roomMsg.ReaderIDs, userData.ID)
			_, err = collection.UpdateOne(context.Background(), bson.M{"_id": roomMsg.ID}, bson.D{{"$set", bson.D{{"reader_ids", readerIDs}}}})
			if err != nil {
				res.Status = "Error"
				res.Message = err.Error()
				json.NewEncoder(c.Writer).Encode(res)
				return
			}
		}

		res.Status = "Success"
		res.Message = "Messages' reader has been added"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func PopulateRoomsFromUserID(id primitive.ObjectID, start int64, limit int64) ([]models.Room, error) {
	var query = bson.M{"user_ids": bson.M{"$elemMatch": bson.M{"$eq": id}}}
	var options = &options.FindOptions{}
	options.SetSort(bson.D{{"last_active", -1}})

	limit = limit - start
	if limit > 24 {
		limit = 24
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

		lastMsg := GetLastRoomMsg(result.ID)
		if lastMsg.Image != "" {
			result.LastMessage = "An image was sent"
		} else {
			result.LastMessage = lastMsg.Message
		}
		result.IsRead = IsLastMessageRead(result.ID, id)
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
		AvatarURL string             `json:"avatar_url" bson:"avatar"`
	}{}

	for _, user := range room.UserIDs {
		err := DB.Collection("users").FindOne(context.Background(), bson.M{"_id": user}).Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		if result.AvatarURL != "" {
			if !strings.HasPrefix(result.AvatarURL, "https://") {
				result.AvatarURL = AvatarURLPrefix + result.ID.Hex() + "/" + result.AvatarURL
			}
		}

		results = append(results, result)
	}

	return results
}

func GetLastRoomMsg(id primitive.ObjectID) models.RoomMessage {
	var roomMsg models.RoomMessage
	var options = &options.FindOneOptions{}
	options.SetSort(bson.D{{"created_at", -1}})

	collection := DB.Collection("room_messages")
	err := collection.FindOne(context.Background(), bson.D{{"room_id", id}}, options).Decode(&roomMsg)
	if err != nil {
		return models.RoomMessage{}
	}

	return roomMsg
}

func IsLastMessageRead(id primitive.ObjectID, userID primitive.ObjectID) bool {
	roomMsg := GetLastRoomMsg(id)

	readerIDs := roomMsg.ReaderIDs
	for _, readerID := range readerIDs {
		if readerID == userID {
			return true
		}
	}
	return false
}

func PopulateRoomMsgFromRoomID(id primitive.ObjectID, start int64, limit int64) []models.RoomMessage {
	var query = bson.M{"room_id": id}

	limit = limit - start

	if start <= 0 {
		start = 0
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

func GetRoomFromUserIDs(firstUserID primitive.ObjectID, secondUserID primitive.ObjectID) (primitive.ObjectID, error) {
	var room models.Room
	collection := DB.Collection("rooms")

	query := bson.D{{"user_ids", bson.D{{"$all", bson.A{firstUserID, secondUserID}}}}}
	roomNotFound := collection.FindOne(context.Background(), query).Decode(&room)
	if roomNotFound != nil {
		room.LastActive = primitive.NewDateTimeFromTime(time.Now())
		room.UserIDs = []primitive.ObjectID{firstUserID, secondUserID}
		result, err := collection.InsertOne(context.TODO(), room)
		if err != nil {
			return primitive.NilObjectID, err
		}

		_, err = DB.Collection("users").UpdateOne(context.Background(), bson.M{"_id": firstUserID}, bson.M{"$set": bson.M{"message_read": false}})
		if err != nil {
			log.Fatal(err)
		}
		_, err = DB.Collection("users").UpdateOne(context.Background(), bson.M{"_id": secondUserID}, bson.M{"$set": bson.M{"message_read": false}})
		if err != nil {
			log.Fatal(err)
		}

		return result.InsertedID.(primitive.ObjectID), nil
	}
	return room.ID, nil
}

func AddMessage(msg models.RoomMessage) (primitive.ObjectID, error) {
	var collection = DB.Collection("room_messages")
	msg.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	msg.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	inserted, err := collection.InsertOne(context.TODO(), msg)

	collection = DB.Collection("rooms")
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": msg.RoomID}, bson.D{{"$set", bson.D{{"last_active", msg.CreatedAt}}}})
	if err != nil {
		return primitive.NilObjectID, err
	}
	notifyUnreadMessage(msg.UserID)
	return inserted.InsertedID.(primitive.ObjectID), nil
}
