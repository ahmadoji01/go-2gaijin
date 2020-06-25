package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RoomMessage struct {
	ID        primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Message   string               `json:"message" bson:"message"`
	Name      string               `json:"name,omitempty" bson:"name,omitempty"`
	Image     string               `json:"image" bson:"image"`
	ImgData   string               `json:"img_data,omitempty"`
	CreatedAt primitive.DateTime   `json:"created_at" bson:"created_at"`
	UserID    primitive.ObjectID   `json:"user_id" bson:"user_id"`
	RoomID    primitive.ObjectID   `json:"room_id" bson:"room_id"`
	ReaderIDs []primitive.ObjectID `json:"reader_ids,omitempty" bson:"reader_ids,omitempty"`
}
