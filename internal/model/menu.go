package model

import "time"

// Menu represents database entity
type Menu struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Calories    int       `json:"calories"`
	Price       float64   `json:"price"`
	Ingredients []string  `gorm:"serializer:json" json:"ingredients"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MenuFilter stores search paramter from query param
type MenuFilter struct {
	Query    string
	Category string
	MinPrice float64
	MaxPrice float64
	MaxCal   int
	Sort     string
	Page     int
	PerPage  int
}

// PaginationResponse helper for output
type PaginationResponse struct {
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}
