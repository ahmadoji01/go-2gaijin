package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductDetail struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	ProductID  primitive.ObjectID `json:"product_id" bson:"product_id"`
	Brand      string             `json:"brand" bson:"brand"`
	Condition  string             `json:"condition" bson:"condition"`
	YearsOwned string             `json:"years_owned" bson:"years_owned"`
	ModelName  string             `json:"model_name" bson:"model_name"`
}
