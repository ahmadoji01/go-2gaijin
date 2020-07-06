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
		subsExpiry := time.Now().AddDate(0, monthsSubscribed, 0)
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
