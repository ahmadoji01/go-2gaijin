package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PaymentMethod struct {
	ID                primitive.ObjectID `json:"_id" bson:"_id"`
	WeChat            string             `json:"wechat" bson:"wechat"`
	COD               bool               `json:"cod" bson:"cod"`
	PayPal            string             `json:"paypal" bson:"paypal"`
	BankAccountNumber string             `json:"bank_account_number" bson:"bank_account_number"`
	BankAccountName   string             `json:"bank_account_name" bson:"bank_account_name"`
	BankName          string             `json:"bank_name" bson:"bank_name"`
	UserID            primitive.ObjectID `json:"user_id" bson:"user_id"`
}
