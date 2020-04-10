package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    primitive.D        `json:"name" bson:"name"`
	IconURL string             `json:"icon_url" bson:"icon_url"`
}
