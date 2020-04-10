package responses

import "gitlab.com/kitalabs/go-2gaijin/models"

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
