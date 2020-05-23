package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Notification struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Name          string             `json:"name" bson:"name"`
	CreatedAt     primitive.DateTime `json:"created_at" bson:"created_at"`
	IsRead        bool               `json:"is_read" bson:"is_read"`
	Type          string             `json:"type" bson:"type"`
	Status        string             `json:"status" bson:"status"`
	NotifIcon     string             `json:"notif_icon,omitempty" bson:"notif_icon,omitempty"`
	NotifiedID    primitive.ObjectID `json:"notified_id" bson:"notified_id"`
	NotifierID    primitive.ObjectID `json:"notifier_id" bson:"notifier_id"`
	AppointmentID primitive.ObjectID `json:"appointment_id,omitempty" bson:"appointment_id,omitempty"`
	Appointment   interface{}        `json:"appointment,omitempty"`
}
