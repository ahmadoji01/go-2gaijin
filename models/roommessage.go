package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoomMessage struct {
	ID        primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Message   string               `json:"message" bson:"message"`
	Name      string               `json:"name,omitempty" bson:"name,omitempty"`
	Image     string               `json:"image,omitempty" bson:"image,omitempty"`
	CreatedAt primitive.DateTime   `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UserID    primitive.ObjectID   `json:"user_id,omitempty" bson:"user_id,omitempty"`
	RoomID    primitive.ObjectID   `json:"room_id,omitempty" bson:"room_id,omitempty"`
	ReaderIDs []primitive.ObjectID `json:"reader_ids,omitempty" bson:"reader_ids,omitempty"`
}
