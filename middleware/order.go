package middleware

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/kitalabs/go-2gaijin/models"
	"gitlab.com/kitalabs/go-2gaijin/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertNotification(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var notification models.Notification
	var res responses.GenericResponse

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &notification)
	if err != nil {
		res.Status = "Error"
		res.Message = "Something wrong happened. Try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	notification.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	notification.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	var collection = DB.Collection("notifications")
	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		notification.NotifierID = userData.ID
		newNotif, err := collection.InsertOne(context.TODO(), notification)

		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		res.Status = "Success"
		res.Message = "Notification successfully saved"
		res.Data = newNotif
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func InsertAppointment(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var appointment models.Appointment
	var res responses.GenericResponse

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &appointment)
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	appointment.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	appointment.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	appointment.ExpiresAt = primitive.NewDateTimeFromTime(time.Now().Add(time.Hour * 24))

	var collection = DB.Collection("appointments")
	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		appointment.RequesterID = userData.ID
		newApp, err := collection.InsertOne(context.TODO(), appointment)

		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		var product models.Product
		collection = DB.Collection("products")
		err = collection.FindOne(context.Background(), bson.M{"_id": appointment.ProductID}).Decode(&product)
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		notifName := userData.FirstName + " wants to meet you at " + appointment.MeetingTime.Time().String() + " for your " + product.Name
		addNotification(notifName, "order_incoming", "", appointment.SellerID, userData.ID, newApp.InsertedID.(primitive.ObjectID))

		res.Status = "Success"
		res.Message = "Appointment successfully saved"
		res.Data = newApp
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func InsertTrustCoin(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var trustcoin models.TrustCoin
	var res responses.GenericResponse

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &trustcoin)
	if err != nil {
		res.Status = "Error"
		res.Message = "Something wrong happened. Try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	trustcoin.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	trustcoin.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	var collection = DB.Collection("trust_coins")
	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		trustcoin.GiverID = userData.ID
		newCoin, err := collection.InsertOne(context.TODO(), trustcoin)

		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		res.Status = "Success"
		res.Message = "Trust coin successfully saved"
		res.Data = newCoin
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func AppointmentConfirmation(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

}

func addNotification(name string, notifType string, notifIcon string, notifiedID primitive.ObjectID, notifierID primitive.ObjectID, appointmentID primitive.ObjectID) {
	var collection = DB.Collection("notifications")
	var notification models.Notification

	notification.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	notification.Name = name
	notification.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	notification.IsRead = false
	notification.Type = notifType
	notification.Status = "pending"
	notification.NotifIcon = notifIcon
	notification.NotifierID = notifierID
	notification.NotifiedID = notifiedID
	notification.AppointmentID = appointmentID

	_, err := collection.InsertOne(context.TODO(), notification)

	if err != nil {
		log.Fatal(err)
	}

	return
}
