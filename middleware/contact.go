package middleware

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/kitalabs/go-2gaijin/config"
	"gitlab.com/kitalabs/go-2gaijin/models"
	"gitlab.com/kitalabs/go-2gaijin/responses"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertTicket(c *gin.Context) {
	c.Writer.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	c.Writer.Header().Set("Access-Control-Allow-Origin", config.CORS)
	c.Writer.Header().Set("Content-Type", "application/json")

	var ticket models.Ticket
	var res responses.GenericResponse

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &ticket)
	if err != nil {
		res.Status = "Error"
		res.Message = err.Error()
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	ticket.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	ticket.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	newTicket, err := DB.Collection("tickets").InsertOne(context.TODO(), ticket)
	if err != nil {
		res.Status = "Error"
		res.Message = "Something wrong happened. Try again"
		json.NewEncoder(c.Writer).Encode(res)
		return
	}

	SendTicketEmail(ticket.Name, ticket.Email, ticket.Message)

	res.Status = "Success"
	res.Message = "Ticket successfully saved"
	res.Data = newTicket
	json.NewEncoder(c.Writer).Encode(res)
	return
}
