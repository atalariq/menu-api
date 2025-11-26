// Package service
package service

import (
	"errors"
	"fmt"
	"log"
	"os"

	"atalariq/menu-api/internal/model"
	"atalariq/menu-api/internal/repository"
	"context"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type MenuService interface {
	Create(input model.Menu) (model.Menu, error)
	GetList(filter model.MenuFilter) (model.PaginationResponse, error)
	GetDetail(id uint) (model.Menu, error)
	Update(id uint, input model.Menu) (model.Menu, error)
	Delete(id uint) error
	GetGrouped(mode string, limit int) (interface{}, error)
	GenerateDescriptionAI(name string, ingredients []string) (string, error)
	GetRecommendationAI(userPreference string) (string, error)
}

type menuService struct {
	repo repository.MenuRepository
}

func NewMenuService(repo repository.MenuRepository) MenuService {
	return &menuService{repo}
}

func (s *menuService) Create(input model.Menu) (model.Menu, error) {
	if input.Price < 0 {
		return model.Menu{}, errors.New("price cannot be negative")
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

	data, pagination, err := s.repo.FindAll(filter)
	pagination.Data = data
	return pagination, err
}

func (s *menuService) GetDetail(id uint) (model.Menu, error) {
	return s.repo.FindByID(id)
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

func (s *menuService) GetGrouped(mode string, limit int) (interface{}, error) {
	return s.repo.GroupBy(mode, limit)
}

// Helper function
func (s *menuService) callGemini(prompt string) (string, error) {
	ctx := context.Background()

	// Get API Keys from env variable
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY not set in environment")
	}

	// Setup client
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Use Gemini flash model
	model := client.GenerativeModel("gemini-2.0-flash")

	// Generate content based on given prompt
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("empty response from AI")
	}

	// Extract and clean up text
	if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		cleanText := strings.TrimSpace(string(textPart))
		cleanText = strings.Trim(cleanText, "\"")

		return cleanText, nil
	}

	return "", errors.New("unexpected response format")
}

func (s *menuService) GenerateDescriptionAI(name string, ingredients []string) (string, error) {
	prompt := fmt.Sprintf(`
		Role: Senior Culinary Copywriter.
		Task: Write a menu description for "%s".

		Ingredients: %s.

		Constraints:
		1. Focus on SENSORY details (texture, temperature, specific flavor notes).
		2. Do NOT use generic words like "delicious", "yummy", or "tasty".
		3. Keep it under 20 words.
		4. Language: English (Elegant & Appetizing).

		Output example: "Silky steamed milk meets robust espresso, finished with a touch of caramelized sweetness."

		Result without any intro or chit-chat:
		`, name, strings.Join(ingredients, ", "))
	return s.callGemini(prompt)
}

func (s *menuService) GetRecommendationAI(userPreference string) (string, error) {
	// Get menu data from datas (context)
	menus, _, err := s.repo.FindAll(model.MenuFilter{PerPage: 100})
	if err != nil {
		return "", err
	}

	var menuListBuilder strings.Builder
	for _, m := range menus {
		menuListBuilder.WriteString(fmt.Sprintf("- Name: %s | Price: %.0f | Category: %s | Ingredients: %s | Desc: %s\n",
			m.Name, m.Price, m.Category, strings.Join(m.Ingredients, ", "), m.Description))
	}

	// Prompting
	prompt := fmt.Sprintf(`
		Role: Strict Menu Recommendation Engine.
    
    Context:
    Menu Data: %s
    User Request: "%s"

    Task: Recommend 1-2 items based on the user request and menu data.
    
    STRICT OUTPUT RULES:
    1. Direct answer ONLY. Do NOT start with "Okay", "Here are", "I understand".
    2. Do NOT use markdown bolding (**).
    3. Format: [Menu Name] - [Benefit in 1 sentence]
    4. Use Menu Data as source of informations. If no match found, output: "No suitable recommendation."
		5. Avoid hallucination, always remember point 4.

    Example Output:
    Es Kopi Susu - The caffeine kick you need to wake up immediately.
    Americano - Zero sugar option for a pure energy boost.

		Result without any intro or chit-chat:
		`, menuListBuilder.String(), userPreference)

	// Return AI Response
	return s.callGemini(prompt)
}
