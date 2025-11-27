// Package service
package service

import (
	"errors"

	"atalariq/menu-api/internal/model"
	"atalariq/menu-api/internal/repository"
)

type MenuService interface {
	Create(input model.Menu) (model.Menu, error)
	GetList(filter model.MenuFilter) (model.PaginationResponse, error)
	GetDetail(id uint) (model.MenuResponse, error)
	Update(id uint, input model.Menu) (model.Menu, error)
	Delete(id uint) error
	GetGrouped(mode string, limit int) (any, error)

	// Add bridge to access `ai_service.go` methods
	GenerateDescription(name string, ingredients []string) (string, error)
	GetRecommendations(request model.RecommendationRequest) ([]model.RecommendationResponse, error)
}

type menuService struct {
	repo repository.MenuRepository
	ai   AIService
}

func NewMenuService(repo repository.MenuRepository, ai AIService) MenuService {
	return &menuService{
		repo: repo,
		ai:   ai,
	}
}

func (s *menuService) Create(input model.Menu) (model.Menu, error) {
	if input.Price < 0 {
		return model.Menu{}, errors.New("price cannot be negative")
	}

	// Use AI to generate description automatically
	if input.Description == "" {
		desc, err := s.ai.GenerateDescription(input.Name, input.Ingredients)
		if err == nil {
			input.Description = desc
		} else { // Fallback if AI throw error
			input.Description = "Delicious " + input.Name
		}
	}
	err := s.repo.Create(&input)
	return input, err
}

func (s *menuService) GetList(filter model.MenuFilter) (model.PaginationResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 {
		filter.PerPage = 10
	}

	menus, pagination, err := s.repo.FindAll(filter)

	if err != nil {
		return model.PaginationResponse{}, err
	}

	var menuResponses []model.MenuResponse
	for _, m := range menus {
		menuResponses = append(menuResponses, m.ToResponse())
	}

	pagination.Data = menuResponses
	return pagination, err
}

func (s *menuService) GetDetail(id uint) (model.MenuResponse, error) {
	menu, err := s.repo.FindByID(id)
	if err != nil {
		return model.MenuResponse{}, err
	}
	// Konversi sebelum return
	return menu.ToResponse(), nil
}

func (s *menuService) Update(id uint, input model.Menu) (model.Menu, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return model.Menu{}, err
	}

	// Update fields
	existing.Name = input.Name
	existing.Price = input.Price
	existing.Calories = input.Calories
	existing.Category = input.Category
	existing.Description = input.Description
	existing.Ingredients = input.Ingredients
	existing.UpdatedAt = input.UpdatedAt

	err = s.repo.Update(&existing)
	return existing, err
}

func (s *menuService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *menuService) GetGrouped(mode string, limit int) (any, error) {
	return s.repo.GroupBy(mode, limit)
}

func (s *menuService) GenerateDescription(name string, ingredients []string) (string, error) {
	return s.ai.GenerateDescription(name, ingredients)
}

func (s *menuService) GetRecommendations(request model.RecommendationRequest) ([]model.RecommendationResponse, error) {
	menus, _, err := s.repo.FindAll(model.MenuFilter{PerPage: 100})
	if err != nil {
		return nil, err
	}

	rawRecommendations, err := s.ai.GetRecommendations(request, menus)
	if err != nil {
		return nil, err
	}

	menuMap := make(map[string]model.Menu)
	for _, m := range menus {
		menuMap[m.Name] = m
	}

	var finalRecommendations []model.RecommendationResponse
	for _, raw := range rawRecommendations {
		if originalMenu, exists := menuMap[raw.MenuName]; exists {
			finalRecommendations = append(finalRecommendations, model.RecommendationResponse{
				Menu:   originalMenu.ToResponse(),
				Reason: raw.Reason,
			})
		}
	}

	return finalRecommendations, nil
}
