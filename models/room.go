package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Room struct {
	ID             primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	RoomType       int                  `json:"room_type,omitempty" bson:"room_type_cd,omitempty"`
	Name           string               `json:"name,omitempty" bson:"name,omitempty"`
	UserIDs        []primitive.ObjectID `json:"user_ids,omitempty" bson:"user_ids,omitempty"`
	LastActive     primitive.DateTime   `json:"last_active,omitempty" bson:"last_active,omitempty"`
	UnreadMessages int                  `json:"unread_messages,omitempty" bson:"unread_messages,omitempty"`
	IsRead         bool                 `json:"is_read" bson:"is_read"`
	IconURL        string               `json:"icon_url"`
	LastMessage    string               `json:"last_message"`
}
