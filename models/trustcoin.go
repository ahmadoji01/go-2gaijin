package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TrustCoin struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Type          string             `json:"type" bson:"type"`
	CreatedAt     primitive.DateTime `json:"created_at" bson:"created_at"`
	GiverID       primitive.ObjectID `json:"giver_id" bson:"giver_id"`
	ReceiverID    primitive.ObjectID `json:"receiver_id" bson:"receiver_id"`
	AppointmentID primitive.ObjectID `json:"appointment_id" bson:"appointment_id"`
}
