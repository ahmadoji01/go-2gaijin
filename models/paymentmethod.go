package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PaymentMethod struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	WeChat      string             `json:"wechat" bson:"wechat"`
	COD         bool               `json:"cod" bson:"cod"`
	PayPal      string             `json:"paypal" bson:"paypal"`
	BankAccount string             `json:"bank_account" bson:"bank_account"`
	BankBranch  string             `json:"bank_branch" bson:"bank_branch"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
}
