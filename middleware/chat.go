package middleware

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
	"gitlab.com/kitalabs/go-2gaijin/channels"
	"gitlab.com/kitalabs/go-2gaijin/config"
	"gitlab.com/kitalabs/go-2gaijin/models"
	"gitlab.com/kitalabs/go-2gaijin/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	_, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		var start, limit int64

		id, err := primitive.ObjectIDFromHex(urlQuery.Get("room"))
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened"
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
			start = limit - 8
			if start < 1 {
				start = 1
			}
		} else {
			start, err = strconv.ParseInt(urlQuery.Get("start"), 10, 64)
		}
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		msgsData = PopulateRoomMsgFromRoomID(id, start, limit)
		roomData.Messages = msgsData
		roomData.TotalMessages = totalMessages
		setIsRead(id)

		res.Status = "Success"
		res.Message = "Chat Messages Successfully Retrieved"
		res.Data = roomData
		json.NewEncoder(c.Writer).Encode(res)
		return
	} else {
		res.Status = "Error"
		res.Message = "Unauthorized"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
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

		res.Status = "Success"
		res.Message = "Message successfully saved"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func InsertPictureMessage(c *gin.Context) {
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

		var roomMsgPic responses.InsertRoomPicture
		roomMsgPic.RoomMsg = roomMsg

		res.Status = "Success"
		res.Message = "Picture Message Successfully Saved"
		res.Data = roomMsgPic
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
	} else {
		res.Status = "Error"
		res.Message = "Unauthorized"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
}
