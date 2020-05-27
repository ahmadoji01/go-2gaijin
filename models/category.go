package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        primitive.D        `json:"name" bson:"name"`
	DisplayName string             `json:"display_name" bson:"display_name"`
	IconURL     string             `json:"icon_url" bson:"icon_url"`
	ParentID    primitive.ObjectID `json:"parent_id" bson:"parent_id"`
	Depth       int64              `json:"depth" bson:"depth"`
}
