package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Delivery struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Origin        string             `json:"origin" bson:"origin"`
	Destination   string             `json:"destination" bson:"destination"`
	Notes         string             `json:"notes" bson:"notes"`
	Name          string             `json:"name" bson:"name"`
	Email         string             `json:"email" bson:"email"`
	Phone         string             `json:"phone" bson:"phone"`
	WeChat        string             `json:"wechat" bson:"wechat"`
	Facebook      string             `json:"facebook" bson:"facebook"`
	DeliveryTime  primitive.DateTime `json:"delivery_time" bson:"delivery_time"`
	CreatedAt     primitive.DateTime `json:"created_at" bson:"created_at"`
	AppointmentID primitive.ObjectID `json:"appointment_id" bson:"appointment_id"`
	RequesterID   primitive.ObjectID `json:"requester_id" bson:"requester_id"`
}

type DeliveryOrder struct {
	Appointment Appointment `json:"appointment"`
	Delivery    Delivery    `json:"delivery"`
}
