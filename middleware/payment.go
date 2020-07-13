package middleware

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"gitlab.com/kitalabs/go-2gaijin/config"
	"gitlab.com/kitalabs/go-2gaijin/models"
	"gitlab.com/kitalabs/go-2gaijin/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreditCardPayment(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
	c.Writer.Header().Set("Content-Type", "application/json")

	var res responses.ResponseMessage
	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		var token responses.OmisePaymentToken
		body, _ := ioutil.ReadAll(c.Request.Body)
		err := json.Unmarshal(body, &token)
		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		client, e := omise.NewClient(config.OmisePublicKey, config.OmiseSecretKey)
		if e != nil {
			log.Fatal(e)
		}

		// Creates a charge from the token
		charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
			Amount:   token.Amount,
			Currency: token.Currency,
			Card:     token.Token,
		}
		if e := client.Do(charge, createCharge); e != nil {
			res.Status = "Error"
			res.Message = e.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		monthsSubscribed := token.MonthsSubscribed
		now := time.Now()
		subsExpiry := now.AddDate(0, monthsSubscribed, 0)
		update := bson.M{"$set": bson.M{"subscription": "basic", "subs_expiry_date": primitive.NewDateTimeFromTime(subsExpiry)}}

		collection := DB.Collection("users")
		_, err = collection.UpdateOne(context.Background(), bson.M{"_id": userData.ID}, update)
		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		res.Status = "Success"
		res.Message = "Payment has been made"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func KonbiniPayment(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
	c.Writer.Header().Set("Content-Type", "application/json")

	var res responses.ResponseMessage
	tokenString := c.Request.Header.Get("Authorization")
	_, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		var source responses.OmisePaymentSource
		body, _ := ioutil.ReadAll(c.Request.Body)
		err := json.Unmarshal(body, &source)
		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		client, e := omise.NewClient(config.OmisePublicKey, config.OmiseSecretKey)
		if e != nil {
			log.Fatal(e)
		}

		// Creates a charge from the token
		charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
			Amount:   source.Amount,
			Currency: source.Currency,
			Source:   source.SourceID,
		}
		if e := client.Do(charge, createCharge); e != nil {
			res.Status = "Error"
			res.Message = e.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		res.Status = "Success"
		res.Message = "Konbini charge has been created"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func KonbiniPaymentSuccessful(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
	c.Writer.Header().Set("Content-Type", "application/json")

	//TODO: Add Konbini Charge Complete Hooks
}

func UpdatePaymentMethod(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
	c.Writer.Header().Set("Content-Type", "application/json")

	var res responses.ResponseMessage
	tokenString := c.Request.Header.Get("Authorization")
	userData, isLoggedIn := LoggedInUser(tokenString)
	if isLoggedIn {
		var paymentMethod models.PaymentMethod
		body, _ := ioutil.ReadAll(c.Request.Body)
		err := json.Unmarshal(body, &paymentMethod)
		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		update := bson.M{"$set": bson.M{"wechat": paymentMethod.WeChat,
			"cod":          paymentMethod.COD,
			"paypal":       paymentMethod.PayPal,
			"bank_account": paymentMethod.BankAccount,
			"bank_branch":  paymentMethod.BankBranch,
		}}

		collection := DB.Collection("payment_methods")
		_, err = collection.UpdateOne(context.Background(), bson.M{"user_id": userData.ID}, update)
		if err != nil {
			res.Status = "Error"
			res.Message = err.Error()
			json.NewEncoder(c.Writer).Encode(res)
			return
		}

		res.Status = "Success"
		res.Message = "Payment method has been updated"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	res.Status = "Error"
	res.Message = "Unauthorized"
	json.NewEncoder(c.Writer).Encode(res)
	return
}

func GetPaymentMethod(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
	c.Writer.Header().Set("Content-Type", "application/json")

	userID, err := primitive.ObjectIDFromHex(c.Param("id"))

	var res responses.GenericResponse
	if err != nil {
		res.Status = "Error"
		res.Message = "Error while converting ID, try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	var paymentMethod models.PaymentMethod
	paymentMethod.UserID = userID
	payment := struct {
		Method models.PaymentMethod `json:"payment_method"`
	}{}

	err = DB.Collection("payment_methods").FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&paymentMethod)
	if err != nil {
		_, err = DB.Collection("payment_methods").InsertOne(context.Background(), paymentMethod)
		if err != nil {
			res.Status = "Error"
			res.Message = "Error while retrieving payment methods, try again"
			json.NewEncoder(c.Writer).Encode(res)
			return
		}
		payment.Method = paymentMethod
		res.Status = "Success"
		res.Message = "Payment method has successfully been retrieved"
		res.Data = payment
		json.NewEncoder(c.Writer).Encode(res)
		return
	}
	payment.Method = paymentMethod

	res.Status = "Success"
	res.Message = "Payment method has successfully been retrieved"
	res.Data = payment
	json.NewEncoder(c.Writer).Encode(res)
	return
}
