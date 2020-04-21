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
	Categories  []primitive.ObjectID `json:"category_ids" bson:"category_ids"`

	DateCreated primitive.DateTime `json:"created_at" bson:"created_at"`
	DateUpdated primitive.DateTime `json:"updated_at" bson:"updated_at"`

	UserID         primitive.ObjectID   `json:"user_id" bson:"user_id"`
	Comments       []primitive.ObjectID `json:"comment_ids" bson:"comment_ids"`
	ProductDetails primitive.ObjectID   `json:"product_details_id" bson:"product_details_id"`
	Keywords       []string             `json:"_keywords" bson:"_keywords"`
	Followers      []primitive.ObjectID `json:"follower_ids" bson:"follower_ids"`
	Orders         []primitive.ObjectID `json:"order_ids" bson:"order_ids"`
	Appointments   []primitive.ObjectID `json:"appointment_ids" bson:"appointment_ids"`

	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`

	PageViews int `json:"page_views"`

	Status ProductStatus `json:"status_cd"`
}

type ProductCard struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Price      int                `json:"price" bson:"price"`
	User       primitive.ObjectID `bson:"user_id"`
	SellerName string             `json:"seller_name"`
	Loc        string             `json:"loc"`
	ImgURL     string             `json:"img_url"`
}
