package models

import (
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    primitive.D        `json:"name" bson:"name"`
	IconURL string             `json:"icon_url" bson:"icon_url"`
}

var CategoryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Category",
		Fields: graphql.Fields{
			"_id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"icon_url": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
