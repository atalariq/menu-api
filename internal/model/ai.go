package model

// RecommendationRequest helper for AI Recommendation
type RecommendationRequest struct {
	Preference string `json:"preference" binding:"required"`
}
