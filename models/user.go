package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email       string             `json:"email" bson:"email"`
	FirstName   string             `json:"first_name" bson:"first_name"`
	LastName    string             `json:"last_name" bson:"last_name"`
	DateOfBirth primitive.DateTime `json:"date_of_birth,omitempty" bson:"date_of_birth,omitempty"`
	Phone       string             `json:"phone" bson:"phone"`
	WeChat      string             `json:"wechat,omitempty" bson:"wechat,omitempty"`
	AvatarURL   string             `json:"avatar_url" bson:"avatar"`
	ShortBio    string             `json:"short_bio" bson:"short_bio"`

	Password     string `json:"password,omitempty" bson:"password,omitempty"`
	Token        string `json:"authentication_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`

	Role string `json:"role,omitempty" bson:"role,omitempty"`

	CreatedAt    primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    primitive.DateTime `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	LastActiveAt primitive.DateTime `json:"last_active_at,omitempty" bson:"last_active_at,omitempty"`

	Products []primitive.ObjectID `json:"products,omitempty" bson:"products,omitempty"`

	Rooms        []primitive.ObjectID `json:"rooms,omitempty" bson:"rooms,omitempty"`
	RoomMessages []primitive.ObjectID `json:"room_messages,omitempty" bson:"room_messages,omitempty"`

	EmailConfirmed bool `json:"email_confirmed" bson:"email_confirmed"`
	PhoneConfirmed bool `json:"phone_confirmed" bson:"phone_confirmed"`

	ConfirmToken       string             `json:"confirm_token,omitempty" bson:"confirm_token,omitempty"`
	ConfirmTokenExpiry primitive.DateTime `json:"confirm_token_expiry,omitempty" bson:"confirm_token_expiry,omitempty"`
	ConfirmSource      string             `json:"confirm_source"`

	Subscription   string             `json:"subscription" bson:"subscription"`
	SubsExpiryDate primitive.DateTime `json:"subs_expiry_date,omitempty" bson:"subs_expiry_date,omitempty"`
	IsSubscribed   bool               `json:"is_subscribed"`

	FollowedProducts []primitive.ObjectID `json:"followed_product_ids,omitempty" bson:"followed_product_ids,omitempty"`
	Notifications    []primitive.ObjectID `json:"notification_ids,omitempty" bson:"notification_ids,omitempty"`
	Appointments     []primitive.ObjectID `json:"appointment_ids,omitempty" bson:"appointment_ids,omitempty"`
	Orders           []primitive.ObjectID `json:"order_ids,omitempty" bson:"order_ids,omitempty"`

	Addresses      []primitive.ObjectID `json:"address_ids,omitempty" bson:"address_ids,omitempty"`
	PrimaryAddress primitive.ObjectID   `json:"primary_address_id,omitempty" bson:"primary_address_id,omitempty"`
	TrustCoins     []primitive.ObjectID `json:"trust_coin_ids,omitempty" bson:"trust_coin_ids,omitempty"`
	GoldCoin       int64                `json:"gold_coin"`
	SilverCoin     int64                `json:"silver_coin"`

	ResetPasswordToken string             `json:"reset_token,omitempty" bson:"reset_password_token,omitempty"`
	ResetTokenExpiry   primitive.DateTime `json:"reset_token_expiry,omitempty" bson:"reset_token_expiry,omitempty"`

	MessageRead bool `json:"message_read" bson:"message_read"`
	NotifRead   bool `json:"notif_read" bson:"notif_read"`

	AuthTokenExpiry    primitive.DateTime `json:"auth_token_expiry,omitempty"`
	RefreshTokenExpiry primitive.DateTime `json:"refresh_token_expiry,omitempty"`

	GoogleSub  string `json:"google_sub,omitempty" bson:"google_sub,omitempty"`
	FacebookID string `json:"facebook_id,omitempty" bson:"facebook_id,omitempty"`
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

type UploadAvatar struct {
	Avatar string `json:"avatar"`
}
