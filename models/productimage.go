package models

import (
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductImage struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Product primitive.ObjectID `json:"product_id" bson:"product_id"`
	ImgURL  string             `json:"image" bson:"image"`
}

var ProductImageType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ProductImage",
		Fields: graphql.Fields{
			"_id": &graphql.Field{
				Type: graphql.String,
			},
			"image": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
