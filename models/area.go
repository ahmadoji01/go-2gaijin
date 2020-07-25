package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Area struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id"`
	Name  string             `json:"name" bson:"name"`
	Label string             `json:"label" bson:"label"`

	Depth             int64              `json:"depth" bson:"depth"`
	ParentID          primitive.ObjectID `json:"parent_id" bson:"parent_id"`
	GeoLocationCenter GeoJson            `json:"geocenter" bson:"geocenter"`
	Range             float64            `json:"range" bson:"range"`
}
