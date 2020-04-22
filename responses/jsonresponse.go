package responses

import "gitlab.com/kitalabs/go-2gaijin/models"

//For Home Page
type HomeData struct {
	Banners          []models.Banner   `json:"banners"`
	Categories       []models.Category `json:"categories"`
	FeaturedItems    []models.Product  `json:"featureditems"`
	RecentItems      []models.Product  `json:"recentitems"`
	RecommendedItems []models.Product  `json:"recommendeditems"`
	FreeItems        []models.Product  `json:"freeitems"`
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
