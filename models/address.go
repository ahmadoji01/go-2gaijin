package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Address struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	Country        string             `json:"country" bson:"country"`
	State          string             `json:"state" bson:"state"`
	City           string             `json:"city" bson:"city"`
	Street         string             `json:"street" bson:"street"`
	PostalCode     string             `json:"postal_code" bson:"postal_code"`
	BuildingNumber string             `json:"building_number" bson:"building_number"`
	Latitude       float64            `json:"latitude" bson:"latitude"`
	Longitude      float64            `json:"longitude" bson:"longitude"`
}
