package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ticket struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Message   string             `json:"message" bson:"message"`
	Email     string             `json:"email" bson:"email"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
}
