package models

import (
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
	ID          primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string               `json:"name" bson:"name"`
	Price       int                  `json:"price" bson:"price"`
	Description string               `json:"description" bson:"description"`
	Category    []primitive.ObjectID `json:"category_ids" bson:"category_ids"`

	DateCreated primitive.DateTime `json:"created_at" bson:"created_at"`
	DateUpdated primitive.DateTime `json:"updated_at" bson:"updated_at"`

	User primitive.ObjectID `json:"user_id" bson:"user_id"`

	Comments       []primitive.ObjectID `json:"comment_ids" bson:"comment_ids"`
	ProductDetails primitive.ObjectID   `json:"product_details_id" bson:"product_details_id"`

	Keywords  []string             `json:"_keywords" bson:"_keywords"`
	Followers []primitive.ObjectID `json:"followers" bson:"followers"`
	LikedBy   []primitive.ObjectID `json:"liked_by" bson:"liked_by"`

	Orders       []primitive.ObjectID `json:"order_ids" bson:"order_ids"`
	Appointments []primitive.ObjectID `json:"appointment_ids" bson:"appointment_ids"`

	Location  []float64 `json:"location" bson:"location"`
	Latitude  float64   `json:"latitude" bson:"latitude"`
	Longitude float64   `json:"longitude" bson:"longitude"`

	PageViews int `json:"page_views"`

	Status       ProductStatus `json:"status" bson:"status"`
	StatusEnum   int           `json:"status_cd" bson:"status_cd"`
	Availability string        `json:"availability" bson:"availability"`
	Relevance    float64
}

type ProductInsert struct {
	Product       Product        `json:"product"`
	ProductImages []ProductImage `json:"product_images"`
	ProductDetail ProductDetail  `json:"product_detail"`
}
