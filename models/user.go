package models

import (
	"github.com/graphql-go/graphql"
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
	Password    string             `json:"password" bson:"password"`
	Token       string             `json:"token" bson:"token"`
	DateOfBirth primitive.DateTime `json:"date_of_birth" bson:"date_of_birth"`
	Phone       string             `json:"phone" bson:"phone"`
	WeChat      string             `json:"wechat" bson:"wechat"`
	CreatedAt   primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt   primitive.DateTime `json:"updated_at" bson:"updated_at"`
	AvatarURL   string             `json:"avatar_url" bson:"avatar_url"`

	Products []Product `json:"products" bson:"products"`

	Rooms        []primitive.ObjectID `json:"rooms" bson:"rooms"`
	RoomMessages []primitive.ObjectID `json:"room_messages" bson:"room_messages"`

	EmailConfirmed bool `json:"email_confirmed" bson:"email_confirmed"`
	PhoneConfirmed bool `json:"phone_confirmed" bson:"phone_confirmed"`

	Subscribed     SubsType           `json:"subscribed" bson:"subscribed"`
	SubsExpiryDate primitive.DateTime `json:"subs_expiry_date" bson:"subs_expiry_date"`

	FollowedProducts []primitive.ObjectID `json:"followed_product_ids" bson:"followed_product_ids"`
	Notifications    []primitive.ObjectID `json:"notification_ids" bson:"notification_ids"`
	Appointments     []primitive.ObjectID `json:"appointment_ids" bson:"appointment_ids"`
	Orders           []primitive.ObjectID `json:"order_ids" bson:"order_ids"`

	Addresses      []primitive.ObjectID `json:"address_ids" bson:"address_ids"`
	PrimaryAddress primitive.ObjectID   `json:"primary_address_id" bson:"primary_address_id"`
	TrustCoins     []primitive.ObjectID `json:"trust_coin_ids" bson:"trust_coin_ids"`
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

var UserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"_id": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"first_name": &graphql.Field{
				Type: graphql.String,
			},
			"last_name": &graphql.Field{
				Type: graphql.String,
			},
			"date_of_birth": &graphql.Field{
				Type: graphql.String,
			},
			"phone": &graphql.Field{
				Type: graphql.String,
			},
			"wechat": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
			},
			"avatar_url": &graphql.Field{
				Type: graphql.String,
			},
			"email_confirmed": &graphql.Field{
				Type: graphql.Boolean,
			},
			"phone_confirmed": &graphql.Field{
				Type: graphql.Boolean,
			},
			"subscribed": &graphql.Field{
				Type: graphql.Int,
			},
			"subs_expiry_date": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
