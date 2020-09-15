package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/kitalabs/go-2gaijin/config"
	"gitlab.com/kitalabs/go-2gaijin/models"
	"gitlab.com/kitalabs/go-2gaijin/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertNotification(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
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

func InsertDelivery(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var delivery models.Delivery
	var res responses.GenericResponse

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &delivery)
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		delivery.ID = primitive.NewObjectIDFromTimestamp(time.Now())
		delivery.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
		delivery.RequesterID = userData.ID

		newDelivery, err := DB.Collection("deliveries").InsertOne(context.TODO(), delivery)
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		localTime, _ := time.LoadLocation("Asia/Tokyo")
		localDeliveryTime := delivery.DeliveryTime.Time().In(localTime).String()
		SendDeliveryRequestEmail("", delivery.Name, delivery.Email, delivery.Phone, delivery.WeChat, delivery.Facebook, delivery.Destination, localDeliveryTime, delivery.Notes)
		SendDeliverySummaryEmail("", delivery.Name, delivery.Email, delivery.Phone, delivery.WeChat, delivery.Facebook, delivery.Destination, localDeliveryTime, delivery.Notes)

		res.Status = "Success"
		res.Message = "Delivery successfully saved"
		res.Data = newDelivery
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func InsertAppointmentWithDelivery(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var appointment models.Appointment
	var delivery models.Delivery
	var deliveryOrder models.DeliveryOrder
	var res responses.GenericResponse

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &deliveryOrder)
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	var collection = DB.Collection("appointments")
	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		appointment = deliveryOrder.Appointment
		delivery = deliveryOrder.Delivery

		appointment.ID = primitive.NewObjectIDFromTimestamp(time.Now())
		appointment.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
		appointment.ExpiresAt = primitive.NewDateTimeFromTime(time.Now().Add(time.Hour * 24))
		appointment.Status = "pending"

		delivery.ID = primitive.NewObjectIDFromTimestamp(time.Now())
		delivery.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
		delivery.DeliveryTime = appointment.MeetingTime
		delivery.AppointmentID = appointment.ID

		appointment.RequesterID = userData.ID
		notifID := primitive.NewObjectIDFromTimestamp(time.Now())
		appointment.NotificationID = notifID
		newApp, err := collection.InsertOne(context.TODO(), appointment)

		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		_, err = DB.Collection("deliveries").InsertOne(context.TODO(), delivery)
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
		addNotification(notifID, notifName, "order_incoming", "", "pending", appointment.SellerID, userData.ID, newApp.InsertedID.(primitive.ObjectID), product.ID)

		var user models.User
		err = DB.Collection("users").FindOne(context.Background(), bson.M{"_id": appointment.SellerID}).Decode(&user)
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		roomID, _ := GetRoomFromUserIDs(appointment.RequesterID, appointment.SellerID)

		//Add Order Message
		var initMsg models.RoomMessage
		initMsg.ID = primitive.NewObjectIDFromTimestamp(time.Now())
		initMsg.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
		initMsg.UserID = userData.ID
		initMsg.RoomID = roomID
		initMsg.Message = "Hi, I am interested in your " + product.Name
		initMsg.Image = ""
		AddMessage(initMsg)

		//Add Product's Image Message
		initMsg.ID = primitive.NewObjectIDFromTimestamp(time.Now())
		initMsg.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
		initMsg.UserID = userData.ID
		initMsg.RoomID = roomID
		initMsg.Message = ""
		initMsg.Image = FindAProductImage(product.ID)
		AddMessage(initMsg)

		localTime, _ := time.LoadLocation("Asia/Tokyo")
		localDeliveryTime := delivery.DeliveryTime.Time().In(localTime).String()
		SendDeliveryRequestEmail(product.Name, delivery.Name, delivery.Email, delivery.Phone, delivery.WeChat, delivery.Facebook, delivery.Destination, localDeliveryTime, delivery.Notes)
		SendBuyingRequestEmail(user.Email, appointment.Source, product.Name)

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

func InsertAppointment(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
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
		notifID := primitive.NewObjectIDFromTimestamp(time.Now())
		appointment.NotificationID = notifID
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

		roomID, _ := GetRoomFromUserIDs(appointment.RequesterID, appointment.SellerID)

		//Add Order Message
		var initMsg models.RoomMessage
		initMsg.ID = primitive.NewObjectIDFromTimestamp(time.Now())
		initMsg.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
		initMsg.UserID = userData.ID
		initMsg.RoomID = roomID
		initMsg.Message = "Hi, I am interested in your " + product.Name
		initMsg.Image = ""
		AddMessage(initMsg)

		//Add Product's Image Message
		initMsg.ID = primitive.NewObjectIDFromTimestamp(time.Now())
		initMsg.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
		initMsg.UserID = userData.ID
		initMsg.RoomID = roomID
		initMsg.Message = ""
		initMsg.Image = FindAProductImage(product.ID)
		AddMessage(initMsg)

		notifName := userData.FirstName + " wants to meet you at " + appointment.MeetingTime.Time().String() + " for your " + product.Name
		addNotification(notifID, notifName, "order_incoming", "", "pending", appointment.SellerID, userData.ID, newApp.InsertedID.(primitive.ObjectID), product.ID)

		var user models.User
		err = DB.Collection("users").FindOne(context.Background(), bson.M{"_id": appointment.SellerID}).Decode(&user)
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}
		SendBuyingRequestEmail(user.Email, appointment.Source, product.Name)

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
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var trustcoin models.TrustCoin
	var res responses.GenericResponse

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &trustcoin)
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
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
		fmt.Println(notif)
		if notif.Type != "give_trust_coin" {
			res.Status = "Error"
			res.Message = "This is not a notification for giving trust coin"
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
				res.Message = err.Error()
				json.NewEncoder(c.Writer).Encode(res)
				return
			}

			var collection = DB.Collection("notifications")
			update := bson.M{"$set": bson.M{"status": "finished"}}
			_, err = collection.UpdateOne(context.Background(), bson.D{{"_id", notif.ID}}, update)
			if err != nil {
				res.Status = "Error"
				res.Message = err.Error()
				json.NewEncoder(c.Writer).Encode(res)
				return
			}

			notifName := FindUserName(trustcoin.GiverID) + " has given you coin"
			addNotification(primitive.NewObjectIDFromTimestamp(time.Now()), notifName, "trust_coin_sent", "", trustcoin.Type, trustcoin.ReceiverID, trustcoin.GiverID, trustcoin.AppointmentID, primitive.NilObjectID)

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
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
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
		update = bson.M{"$set": bson.M{"status": status}}
		_, err := collection.UpdateOne(context.Background(), bson.D{{"_id", appointment.NotificationID}}, update)
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

		var requester models.User
		collection = DB.Collection("users")
		err = collection.FindOne(context.Background(), bson.M{"_id": appointment.RequesterID}).Decode(&requester)
		if err != nil {
			res.Status = "Error"
			res.Message = "Something wrong happened. Try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		if status == "accepted" {
			setNotifsToRejectOrder(appointment.NotificationID, appointment.ProductID)
			SendAcceptedOrderEmail(product.Name, requester.Email)
		} else if status == "rejected" {
			notifName := "Appointment Rejected"
			addNotification(primitive.NewObjectIDFromTimestamp(time.Now()), notifName, "appointment_confirmation", "", "rejected", appointment.RequesterID, userData.ID, appointment.ID, appointment.ProductID)
			SendAcceptedOrderEmail(product.Name, requester.Email)
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
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
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
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
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
		addNotification(primitive.NewObjectIDFromTimestamp(time.Now()), notifName, "give_trust_coin", "", "pending", appointment.RequesterID, userData.ID, appointment.ID, appointment.ProductID)

		notifName = FindUserName(appointment.SellerID) + " has finished transaction with you"
		addNotification(primitive.NewObjectIDFromTimestamp(time.Now()), notifName, "give_trust_coin", "", "pending", appointment.SellerID, appointment.RequesterID, appointment.ID, appointment.ProductID)

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

func trustCoinNotif(notifiedID primitive.ObjectID, notifierID primitive.ObjectID, appointmentID primitive.ObjectID) models.Notification {
	var collection = DB.Collection("notifications")

	var notif models.Notification
	collection.FindOne(context.Background(), bson.D{{"type", "give_trust_coin"}, {"notified_id", notifiedID}, {"notifier_id", notifierID}, {"appointment_id", appointmentID}}).Decode(&notif)
	return notif
}

func setNotifsToRejectOrder(acceptedNotifID primitive.ObjectID, productID primitive.ObjectID) {
	var collection = DB.Collection("notifications")
	cur, err := collection.Find(context.Background(), bson.M{"product_id": productID})
	if err != nil {
		log.Fatal(err)
	}

	collection = DB.Collection("appointments")
	for cur.Next(context.Background()) {
		var result models.Notification
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}

		if result.ID != acceptedNotifID {
			if result.Status == "pending" && result.Type == "order_incoming" {
				update := bson.M{"$set": bson.M{"status": "rejected"}}
				_, e = collection.UpdateOne(context.Background(), bson.M{"_id": result.AppointmentID}, update)
				_, e = DB.Collection("notifications").UpdateOne(context.Background(), bson.M{"_id": result.ID}, update)
				notifName := "Appointment Rejected"
				addNotification(primitive.NewObjectIDFromTimestamp(time.Now()), notifName, "appointment_confirmation", "", "rejected", result.NotifierID, result.NotifiedID, result.AppointmentID, result.ProductID)
			}
		} else {
			update := bson.M{"$set": bson.M{"status": "accepted"}}
			_, e = collection.UpdateOne(context.Background(), bson.M{"_id": result.AppointmentID}, update)
			_, e = DB.Collection("notifications").UpdateOne(context.Background(), bson.M{"_id": result.ID}, update)
			notifName := "Appointment Accepted"
			addNotification(primitive.NewObjectIDFromTimestamp(time.Now()), notifName, "appointment_confirmation", "", "accepted", result.NotifierID, result.NotifiedID, result.AppointmentID, result.ProductID)
		}
	}
}

func addNotification(notifID primitive.ObjectID, name string, notifType string, notifIcon string, status string, notifiedID primitive.ObjectID, notifierID primitive.ObjectID, appointmentID primitive.ObjectID, productID primitive.ObjectID) {
	var collection = DB.Collection("notifications")
	var notification models.Notification

	notification.ID = notifID
	notification.Name = name
	notification.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	notification.IsRead = false
	notification.Type = notifType
	notification.Status = status
	notification.NotifIcon = notifIcon
	notification.NotifierID = notifierID
	notification.NotifiedID = notifiedID
	notification.AppointmentID = appointmentID
	notification.ProductID = productID

	_, err := collection.InsertOne(context.TODO(), notification)

	_, err = DB.Collection("users").UpdateOne(context.Background(), bson.M{"_id": notification.NotifiedID}, bson.M{"$set": bson.M{"notif_read": false}})
	if err != nil {
		log.Fatal(err)
	}

	return
}
