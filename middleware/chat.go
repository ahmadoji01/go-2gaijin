package middleware

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/kitalabs/go-2gaijin/models"
	"gitlab.com/kitalabs/go-2gaijin/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetChatRoomMsg(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", CORS)
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

		if urlQuery.Get("limit") == "" {
			limit, _ = collection.CountDocuments(context.Background(), bson.D{{"room_id", id}})
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

func InsertMessage(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var roomMsg models.RoomMessage
	var res responses.GenericResponse

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &roomMsg)
	if err != nil {
		res.Status = "Error"
		res.Message = "Something wrong happened. Try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	roomMsg.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	var collection = DB.Collection("room_messages")
	tokenString := c.Request.Header.Get("Authorization")
	_, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		_, err := collection.InsertOne(context.TODO(), roomMsg)
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		res.Status = "Success"
		res.Message = "Message successfully saved"
		json.NewEncoder(c.Writer).Encode(res)
		return
	} else {
		res.Status = "Error"
		res.Message = "Unauthorized"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
}
