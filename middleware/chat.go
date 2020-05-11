package middleware

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/kitalabs/go-2gaijin/models"
	"gitlab.com/kitalabs/go-2gaijin/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetChatRoomMsg(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
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
