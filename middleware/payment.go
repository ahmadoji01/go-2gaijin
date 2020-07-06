package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"gitlab.com/kitalabs/go-2gaijin/config"
	"gitlab.com/kitalabs/go-2gaijin/responses"
)

func CreditCardPayment(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
	c.Writer.Header().Set("Content-Type", "application/json")

	var res responses.ResponseMessage
	tokenString := c.Request.Header.Get("Authorization")
	_, isLoggedIn := LoggedInUser(tokenString)
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

		fmt.Println(token.Amount)

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
