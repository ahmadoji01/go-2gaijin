package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Banner struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ImgURL    string             `json:"img_url" bson:"img_url"`
	BannerURL string             `json:"banner_url" bson:"banner_url"`
}
