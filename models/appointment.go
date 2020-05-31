package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Appointment struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Status     string             `json:"status" bson:"status"`
	CreatedAt  primitive.DateTime `json:"created_at" bson:"created_at"`
	ExpiresAt  primitive.DateTime `json:"expires_at" bson:"expires_at"`
	IsDelivery bool               `json:"is_delivery" bson:"is_delivery"`

	MeetingTime primitive.DateTime `json:"meeting_time" bson:"meeting_time"`
	MeetingLat  float64            `json:"latitude,omitempty" bson:"latitude,omitempty"`
	MeetingLng  float64            `json:"longitude,omitempty" bson:"longitude,omitempty"`

	ProductID       primitive.ObjectID `json:"product_id" bson:"product_id"`
	ProductDetail   interface{}        `json:"product_detail"`
	SellerID        primitive.ObjectID `json:"seller_id" bson:"seller_id"`
	AppointmentUser interface{}        `json:"appointment_user,omitempty"`
	RequesterID     primitive.ObjectID `json:"requester_id" bson:"requester_id"`
}
