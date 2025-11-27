package service

import (
	"atalariq/menu-api/internal/model"
)

type AIService interface {
	GenerateDescription(name string, ingredients []string) (string, error)
	GetRecommendations(request model.RecommendationRequest, menus []model.Menu) ([]model.RecommendationResponseRaw, error)
}
