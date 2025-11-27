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

// MenuResponse used for the API response
type MenuResponse struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Calories    int      `json:"calories"`
	Price       float64  `json:"price"`
	Ingredients []string `json:"ingredients"`
	Description string   `json:"description"`
}

// Helper method to convert Model to Response
func (m *Menu) ToResponse() MenuResponse {
	return MenuResponse{
		ID:          m.ID,
		Name:        m.Name,
		Category:    m.Category,
		Calories:    m.Calories,
		Price:       m.Price,
		Ingredients: m.Ingredients,
		Description: m.Description,
	}
}

type MenuSuccessResponse struct {
	Message string `json:"message"`
	Data    Menu   `json:"data"`
}

type MenuDetailResponse struct {
	Data MenuResponse `json:"data"`
}

type MenuQueryRequest struct {
	Q        string  `form:"q"`
	Category string  `form:"category"`
	MinPrice float64 `form:"min_price"`
	MaxPrice float64 `form:"max_price"`
	MaxCal   int     `form:"max_cal"`
	Sort     string  `form:"sort"`
	Page     int     `form:"page,default=1"`
	PerPage  int     `form:"per_page,default=10"`
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

// MenuPaginationResponse helper for output
type MenuPaginationResponse struct {
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	TotalPages int            `json:"total_pages"`
	Data       []MenuResponse `json:"data"`
}
