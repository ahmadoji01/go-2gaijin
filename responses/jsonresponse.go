package responses

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

type ProductsData struct {
	Products interface{}
}
