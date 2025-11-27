package service

import (
	"atalariq/menu-api/internal/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type geminiService struct {
	apiKey string
}

func NewGeminiService() AIService {
	return &geminiService{
		apiKey: os.Getenv("GEMINI_API_KEY"),
	}
}

func (s *geminiService) callGemini(prompt string) (string, error) {
	ctx := context.Background()

	// Setup client
	client, err := genai.NewClient(ctx, option.WithAPIKey(s.apiKey))
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

func (s *geminiService) GenerateDescription(name string, ingredients []string) (string, error) {
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

func (s *geminiService) GetRecommendations(request model.RecommendationRequest, menus []model.Menu) ([]model.RecommendationResponseRaw, error) {
	menuMap := make(map[string]model.Menu)
	var menuListBuilder strings.Builder

	for _, m := range menus {
		menuMap[m.Name] = m

		menuListBuilder.WriteString(fmt.Sprintf("- %s (Ingredients: %s, Category: %s)\n",
			m.Name, strings.Join(m.Ingredients, ", "), m.Category))
	}

	userPreference := request.Preference
	prompt := fmt.Sprintf(`
	Role: Strict Menu Recommendation Engine.
	Context:
	User Request: "%s"
	Available Menu:
	%s

	Task: Recommend 1-3 items based on the user request.

	CRITICAL INSTRUCTION:
	1. Output MUST be a valid JSON Array.
	2. Use the EXACT menu name from the list above.
	3. Format: [{"menu_name": "Exact Name", "reason": "Why it fits"}]
	4. No Markdown. No Intro.
	`, userPreference, menuListBuilder.String())

	rawResponse, err := s.callGemini(prompt)
	if err != nil {
		return nil, err
	}

	cleanJSON := strings.TrimSpace(rawResponse)
	cleanJSON = strings.TrimPrefix(cleanJSON, "```json")
	cleanJSON = strings.TrimPrefix(cleanJSON, "```")
	cleanJSON = strings.TrimSuffix(cleanJSON, "```")

	var rawRecommendations []model.RecommendationResponseRaw
	if err := json.Unmarshal([]byte(cleanJSON), &rawRecommendations); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %v", err)
	}

	return rawRecommendations, nil
}
