package models

import (
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductStatus int

const (
	Available ProductStatus = iota
	Hold
	Sold
)

func (ps ProductStatus) String() string {
	return [...]string{"Available", "Hold", "Sold"}[ps]
}

type Product struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Price       int                `json:"price" bson:"price"`
	Description string             `json:"description" bson:"description"`
	Category    Category           `json:"category" bson:"category"`

	DateCreated primitive.DateTime `json:"created_at" bson:"created_at"`
	DateUpdated primitive.DateTime `json:"updated_at" bson:"updated_at"`

	ProductImages []ProductImage `json:"product_images" bson:"product_images"`

	User User `json:"user" bson:"user"`

	Comments       []primitive.ObjectID `json:"comment_ids" bson:"comment_ids"`
	ProductDetails primitive.ObjectID   `json:"product_details_id" bson:"product_details_id"`

	Keywords  []string `json:"_keywords" bson:"_keywords"`
	Followers []User   `json:"followers" bson:"followers"`

	Orders       []primitive.ObjectID `json:"order_ids" bson:"order_ids"`
	Appointments []primitive.ObjectID `json:"appointment_ids" bson:"appointment_ids"`

	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Location  string `json:"location" bson:"location"`

	PageViews int `json:"page_views"`

	Status ProductStatus `json:"status_cd"`
}

type ProductCard struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name" bson:"name"`
	Price  int                `json:"price" bson:"price"`
	Seller User               `json:"seller"`
	Loc    string             `json:"loc"`
	ImgURL string             `json:"img_url"`
}

var ProductType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Product",
		Fields: graphql.Fields{
			"_id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"price": &graphql.Field{
				Type: graphql.Int,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"category": &graphql.Field{
				Type: CategoryType,
			},
			"product_images": &graphql.Field{
				Type: graphql.NewList(ProductImageType),
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
			},
			"user": &graphql.Field{
				Type: UserType,
			},
			"_keywords": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"status_cd": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
