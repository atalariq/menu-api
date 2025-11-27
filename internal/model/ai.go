package model

// RecommendationRequest stores parameter for AI recommendation request
type RecommendationRequest struct {
	Preference string `json:"preference" binding:"required"`
}

// RecommendationResponseRaw is a helper to catch AI response (token saving and more accurate)
type RecommendationResponseRaw struct {
	MenuName string `json:"menu_name"`
	Reason   string `json:"reason"`
}

// MenuResponse used for the AI recommendation response
type RecommendationResponse struct {
	Menu   MenuResponse `json:"menu"`
	Reason string       `json:"reason"`
}

type RecommendationListResponse struct {
	Data []RecommendationResponse `json:"data"`
}

type GenerateDescriptionRequest struct {
	Name        string   `json:"name" binding:"required"`
	Ingredients []string `json:"ingredients" binding:"required"`
}

type GenerateDescriptionResponse struct {
	Description string `json:"generated_description"`
}
