package model

// dto.go
// Description:
// Provides DTO/struct for Swagger to generate API Documentation automatically

type GeneralResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid input format or ID not found"`
}
