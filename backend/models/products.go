package models

import "time"

type ProductStatus int

const (
	Available ProductStatus = iota
	Hold
	Sold
)

func (ps ProductStatus) String() string {
	return [...]string{"Available", "Hold", "Sold"}[ps]
}

type Product struct {
	Name        string    `json:"name"`
	Price       int       `json:"price"`
	Description string    `json:"description"`
	DateCreated time.Time `json:"created_at"`
	DateUpdated time.Time `json:"updated_at"`

	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`

	PageViews int `json:"page_views"`

	Status ProductStatus `json:"status"`
}
