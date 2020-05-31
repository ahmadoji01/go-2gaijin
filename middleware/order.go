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
	appointment.Status = "pending"

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
		addNotification(notifName, "order_incoming", "", "pending", appointment.SellerID, userData.ID, newApp.InsertedID.(primitive.ObjectID))

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
		var notif = trustCoinNotif(userData.ID, trustcoin.ReceiverID, trustcoin.AppointmentID)

		if notif.Type != "give_trust_coin" {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		if notif.Status == "finished" {
			res.Status = "Error"
			res.Message = "You have given your trust coin for this transaction"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		if notif.Status == "pending" {
			trustcoin.GiverID = userData.ID
			newCoin, err := collection.InsertOne(context.TODO(), trustcoin)

			if err != nil {
				res.Status = "Error"
				res.Message = "Something wrong happened. Try again"
				json.NewEncoder(c.Writer).Encode(res)
				return
			}

			var collection = DB.Collection("notifications")
			update := bson.M{"$set": bson.M{"status": "finished"}}
			_, err = collection.UpdateOne(context.Background(), bson.D{{"_id", notif.ID}}, update)
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

	var appointment models.Appointment
	var res responses.GenericResponse

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &appointment)
	if err != nil {
		res.Status = "Error"
		res.Message = "Something wrong happened. Try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	status := appointment.Status

	var collection = DB.Collection("appointments")
	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		err = collection.FindOne(context.Background(), bson.D{{"_id", appointment.ID}}).Decode(&appointment)
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		if userData.ID != appointment.SellerID {
			res.Status = "Error"
			res.Message = "You are not authorized to confirm this appointment"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		if appointment.Status != "pending" {
			res.Status = "Error"
			res.Message = "You have already confirmed this appointment"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		update := bson.M{"$set": bson.M{"status": status}}
		_, err = collection.UpdateOne(context.Background(), bson.D{{"_id", appointment.ID}}, update)
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		collection = DB.Collection("notifications")
		update = bson.M{"$set": bson.M{"status": appointment.Status}}
		_, err = collection.UpdateOne(context.Background(), bson.D{{"appointment_id", appointment.ID}}, update)
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

		if status == "accepted" {
			notifName := userData.FirstName + " accepted your appointment request for " + product.Name
			addNotification(notifName, "appointment_confirmation", "", "accepted", appointment.RequesterID, userData.ID, appointment.ID)
		} else if status == "rejected" {
			notifName := userData.FirstName + " rejected your appointment request for " + product.Name
			addNotification(notifName, "appointment_confirmation", "", "rejected", appointment.RequesterID, userData.ID, appointment.ID)
		}

		res.Status = "Success"
		res.Message = "Appointment successfully updated"
		res.Data = appointment
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func RescheduleAppointment(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var appointment models.Appointment
	var res responses.GenericResponse

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &appointment)
	if err != nil {
		res.Status = "Error"
		res.Message = "Something wrong happened. Try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	var collection = DB.Collection("appointments")
	tokenString := c.Request.Header.Get("Authorization")
	_, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		update := bson.M{"$set": bson.M{"meeting_time": appointment.MeetingTime, "status": "accepted"}}
		_, err = collection.UpdateOne(context.Background(), bson.D{{"_id", appointment.ID}}, update)
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		var collection = DB.Collection("notifications")
		update = bson.M{"$set": bson.M{"status": "accepted"}}
		_, err = collection.UpdateOne(context.Background(), bson.D{{"appointment_id", appointment.ID}}, update)
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		res.Status = "Success"
		res.Message = "Appointment successfully rescheduled"
		res.Data = appointment
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func FinishAppointment(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var appointment models.Appointment
	var res responses.GenericResponse

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &appointment)
	if err != nil {
		res.Status = "Error"
		res.Message = "Something wrong happened. Try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	var collection = DB.Collection("appointments")
	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		err = collection.FindOne(context.Background(), bson.D{{"_id", appointment.ID}}).Decode(&appointment)
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		if userData.ID != appointment.SellerID {
			res.Status = "Error"
			res.Message = "You are not authorized to finish this appointment"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		if appointment.Status != "accepted" {
			res.Status = "Error"
			res.Message = "You cannot finish this appointment"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		update := bson.M{"$set": bson.M{"status": "finished"}}
		_, err = collection.UpdateOne(context.Background(), bson.D{{"_id", appointment.ID}}, update)
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		collection = DB.Collection("notifications")
		update = bson.M{"$set": bson.M{"status": appointment.Status}}
		_, err = collection.UpdateOne(context.Background(), bson.D{{"appointment_id", appointment.ID}}, update)
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

		notifName := userData.FirstName + " has finished transaction with you"
		addNotification(notifName, "give_trust_coin", "", "pending", appointment.RequesterID, userData.ID, appointment.ID)

		notifName = FindUserName(appointment.SellerID) + " has finished transaction with you"
		addNotification(notifName, "give_trust_coin", "", "pending", appointment.SellerID, appointment.RequesterID, appointment.ID)

		res.Status = "Success"
		res.Message = "Appointment successfully updated"
		res.Data = appointment
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func trustCoinNotif(notifierID primitive.ObjectID, notifiedID primitive.ObjectID, appointmentID primitive.ObjectID) models.Notification {
	var collection = DB.Collection("notifications")

	var notif models.Notification
	collection.FindOne(context.Background(), bson.D{{"type", "give_trust_coin"}, {"notified_id", notifiedID}, {"notifier_id", notifierID}, {"appointment_id", appointmentID}}).Decode(&notif)
	return notif
}

func addNotification(name string, notifType string, notifIcon string, status string, notifiedID primitive.ObjectID, notifierID primitive.ObjectID, appointmentID primitive.ObjectID) {
	var collection = DB.Collection("notifications")
	var notification models.Notification

	notification.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	notification.Name = name
	notification.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	notification.IsRead = false
	notification.Type = notifType
	notification.Status = status
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
