package responses

import (
	"gitlab.com/kitalabs/go-2gaijin/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//For Home Page
type HomeData struct {
	Banners          interface{} `json:"banners"`
	Categories       interface{} `json:"categories"`
	FeaturedItems    interface{} `json:"featureditems"`
	RecentItems      interface{} `json:"recentitems"`
	RecommendedItems interface{} `json:"recommendeditems"`
	FreeItems        interface{} `json:"freeitems"`
}

type HomePage struct {
	Status  string `json:"status"`
	Message string `json:"message"`

	Data HomeData `json:"data"`
}

// For Search Page
type Pagination struct {
	CurrentPage  int64 `json:"current_page"`
	NextPage     int64 `json:"next_page"`
	PreviousPage int64 `json:"previous_page"`
	TotalPages   int64 `json:"total_pages"`
	ItemsPerPage int64 `json:"items_per_page"`
	TotalItems   int64 `json:"total_items"`
}

type SearchPage struct {
	Status  string `json:"status"`
	Message string `json:"message"`

	Data interface{} `json:"data"`
}

type UserData struct {
	User interface{} `json:"user"`
}

type LoginPage struct {
	Status  string `json:"status"`
	Message string `json:"message"`

	Data UserData `json:"data"`
}

type ResponseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ProductCategory struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	Name    string             `json:"name" bson:"name"`
	IconURL string             `json:"icon_url" bson:"icon_url"`
}

type ProductDetailItem struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       int                `json:"price" bson:"price"`
	Category    ProductCategory    `json:"category"`

	Images []interface{} `json:"images"`

	DateCreated primitive.DateTime `json:"created_at" bson:"created_at"`
	DateUpdated primitive.DateTime `json:"updated_at" bson:"updated_at"`

	User primitive.ObjectID `json:"user_id" bson:"user_id"`

	Latitude  string      `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude string      `json:"longitude,omitempty" bson:"longitude,omitempty"`
	Location  interface{} `json:"location"`

	Comments   []interface{} `json:"comment_ids" bson:"comment_ids"`
	StatusEnum int           `json:"status_enum" bson:"status_cd"`
	Status     string        `json:"status"`
}

type ProductDetails struct {
	Item         ProductDetailItem `json:"item"`
	RelatedItems []interface{}     `json:"relateditems"`
	SellerItems  []interface{}     `json:"selleritems"`
}

type ProductDetailPage struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ProfileData struct {
	Profile     models.User `json:"profile"`
	PostedItems interface{} `json:"posted_items"`
}

type SearchData struct {
	Items interface{} `json:"items"`
}

type AppointmentData struct {
	Appointments []models.Appointment `json:"appointments"`
}

type NotificationData struct {
	Notifications []models.Notification `json:"notifications"`
}

type ChatLobbyData struct {
	ChatLobby []models.Room `json:"chat_lobby"`
}

type ChatRoomData struct {
	TotalMessages int64                `json:"total_messages"`
	Messages      []models.RoomMessage `json:"messages"`
}

type RoomMsgData struct {
	Message models.RoomMessage `json:"message"`
}

type GenericResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
