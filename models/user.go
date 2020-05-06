package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubsType int

const (
	Normal SubsType = iota
	Basic
	Full
)

func (st SubsType) String() string {
	return [...]string{"Normal", "Basic", "Full"}[st]
}

type User struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email       string             `json:"email" bson:"email"`
	FirstName   string             `json:"first_name" bson:"first_name"`
	LastName    string             `json:"last_name" bson:"last_name"`
	Password    string             `json:"password,omitempty" bson:"password,omitempty"`
	Token       string             `json:"authentication_token,omitempty" bson:"token,omitempty"`
	DateOfBirth primitive.DateTime `json:"date_of_birth,omitempty" bson:"date_of_birth,omitempty"`
	Phone       string             `json:"phone,omitempty" bson:"phone,omitempty"`
	WeChat      string             `json:"wechat,omitempty" bson:"wechat,omitempty"`
	CreatedAt   primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   primitive.DateTime `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	AvatarURL   string             `json:"avatar_url,omitempty" bson:"avatar,omitempty"`

	Products []primitive.ObjectID `json:"products,omitempty" bson:"products,omitempty"`

	Rooms        []primitive.ObjectID `json:"rooms,omitempty" bson:"rooms,omitempty"`
	RoomMessages []primitive.ObjectID `json:"room_messages,omitempty" bson:"room_messages,omitempty"`

	EmailConfirmed bool `json:"email_confirmed,omitempty" bson:"email_confirmed,omitempty"`
	PhoneConfirmed bool `json:"phone_confirmed,omitempty" bson:"phone_confirmed,omitempty"`

	Subscribed     SubsType           `json:"subscribed,omitempty" bson:"subscribed,omitempty"`
	SubsExpiryDate primitive.DateTime `json:"subs_expiry_date,omitempty" bson:"subs_expiry_date,omitempty"`

	FollowedProducts []primitive.ObjectID `json:"followed_product_ids,omitempty" bson:"followed_product_ids,omitempty"`
	Notifications    []primitive.ObjectID `json:"notification_ids,omitempty" bson:"notification_ids,omitempty"`
	Appointments     []primitive.ObjectID `json:"appointment_ids,omitempty" bson:"appointment_ids,omitempty"`
	Orders           []primitive.ObjectID `json:"order_ids,omitempty" bson:"order_ids,omitempty"`

	Addresses      []primitive.ObjectID `json:"address_ids,omitempty" bson:"address_ids,omitempty"`
	PrimaryAddress primitive.ObjectID   `json:"primary_address_id,omitempty" bson:"primary_address_id,omitempty"`
	TrustCoins     []primitive.ObjectID `json:"trust_coin_ids,omitempty" bson:"trust_coin_ids,omitempty"`

	ResetPasswordToken string             `json:"reset_token,omitempty" bson:"reset_password_token,omitempty"`
	ResetTokenExpiry   primitive.DateTime `json:"reset_token_expiry,omitempty" bson:"reset_token_expiry,omitempty"`
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}
