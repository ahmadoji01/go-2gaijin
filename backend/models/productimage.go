package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductImage struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Product primitive.ObjectID `json:"product_id" bson:"product_id"`
	ImgURL  string             `json:"image" bson:"image"`
}
