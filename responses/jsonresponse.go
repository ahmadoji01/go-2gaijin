package responses

import "go.mongodb.org/mongo-driver/bson/primitive"

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

type SearchData struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

type SearchPage struct {
	Status  string `json:"status"`
	Message string `json:"message"`

	Data SearchData `json:"data"`
}

type UserData struct {
	User interface{} `json:"user"`
}

type LoginPage struct {
	Status  string `json:"status"`
	Message string `json:"message"`

	Data UserData `json:"data"`
}

type ProductDetailPage struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description"`
	Price       int                `json:"price"`
	Category    interface{}        `json:"category"`

	DateCreated primitive.DateTime `json:"created_at" bson:"created_at"`
	DateUpdated primitive.DateTime `json:"updated_at" bson:"updated_at"`

	User primitive.ObjectID `json:"user_id" bson:"user_id"`

	Comments       []interface{}      `json:"comment_ids" bson:"comment_ids"`
	ProductDetails primitive.ObjectID `json:"product_details_id" bson:"product_details_id"`

	Location  []float64 `json:"location" bson:"location"`
	PageViews int       `json:"page_views"`

	Status int `json:"status_enum" bson:"status_cd"`
}
