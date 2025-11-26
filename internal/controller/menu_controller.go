// Package controller
package controller

import (
	"net/http"
	"strconv"

	"atalariq/menu-api/internal/model"
	"atalariq/menu-api/internal/service"

	"github.com/gin-gonic/gin"
)

type MenuController struct {
	service service.MenuService
}

func NewMenuController(service service.MenuService) *MenuController {
	return &MenuController{service}
}

// Create godoc
// @Summary      Create a new menu
// @Description  Create a new menu item with ingredients
// @Tags         menu
// @Accept       json
// @Produce      json
// @Param        menu body model.Menu true "Menu Request"
// @Success      201  {object}  model.Menu
// @Failure      400  {object}  map[string]string
// @Router       /menu [post]
func (c *MenuController) Create(ctx *gin.Context) {
	var input model.Menu
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.service.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Menu created successfully",
		"data":    result,
	})
}

// FindAll godoc
// @Summary      List all menus
// @Description  Get menu list with pagination, filtering, and sorting
// @Tags         menu
// @Produce      json
// @Param        q query string false "Search by name or description"
// @Param        category query string false "Filter by category"
// @Param        min_price query number false "Minimum price"
// @Param        max_price query number false "Maximum price"
// @Param        max_cal query int false "Maximum calories"
// @Param        sort query string false "Sort by field (e.g., price:asc)"
// @Param        page query int false "Page number (default 1)"
// @Param        per_page query int false "Items per page (default 10)"
// @Success      200  {object}  model.PaginationResponse
// @Router       /menu [get]
func (c *MenuController) FindAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	minPrice, _ := strconv.ParseFloat(ctx.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(ctx.Query("max_price"), 64)
	maxCal, _ := strconv.Atoi(ctx.Query("max_cal"))

	filter := model.MenuFilter{
		Query:    ctx.Query("q"),
		Category: ctx.Query("category"),
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		MaxCal:   maxCal,
		Sort:     ctx.Query("sort"),
		Page:     page,
		PerPage:  perPage,
	}

	result, err := c.service.GetList(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Langsung return struct PaginationResponse agar format JSON sesuai standar
	ctx.JSON(http.StatusOK, result)
}

// GetByID godoc
// @Summary      Get menu detail
// @Description  Get details of a specific menu item by ID
// @Tags         menu
// @Produce      json
// @Param        id   path      int  true  "Menu ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]string
// @Router       /menu/{id} [get]
func (c *MenuController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	menu, err := c.service.GetDetail(uint(id))
	if err != nil {
		// Asumsi service return error jika not found
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Menu not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": menu})
}

// Update godoc
// @Summary      Update menu
// @Description  Update an existing menu item
// @Tags         menu
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Menu ID"
// @Param        menu body model.Menu true "Update Data"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]string
// @Router       /menu/{id} [put]
func (c *MenuController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var input model.Menu
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedMenu, err := c.service.Update(uint(id), input)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Menu not found or update failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Menu updated successfully",
		"data":    updatedMenu,
	})
}

// Delete godoc
// @Summary      Delete menu
// @Description  Delete a menu item by ID
// @Tags         menu
// @Produce      json
// @Param        id   path      int  true  "Menu ID"
// @Success      200  {object}  map[string]string
// @Router       /menu/{id} [delete]
func (c *MenuController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Menu not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Menu deleted successfully"})
}

// GroupByCategory godoc
// @Summary      Group menus by category
// @Description  Get menu counts or lists grouped by category
// @Tags         menu
// @Produce      json
// @Param        mode query string true "Mode: 'count' or 'list'"
// @Param        limit query int false "Limit item per category (default 5)"
// @Success      200  {object}  map[string]interface{}
// @Router       /menu/group-by-category [get]
func (c *MenuController) GroupByCategory(ctx *gin.Context) {
	limit, err := strconv.Atoi(ctx.Query("per_category"))
	if err != nil {
		limit = 5
	}

	mode := ctx.Query("mode")
	if mode != "count" && mode != "list" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mode. Use 'count' or 'list'"})
		return
	}

	result, err := c.service.GetGrouped(mode, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// GenerateDescriptionAI godoc
// @Summary      Generate Menu Description (AI)
// @Description  Use Gemini AI to create a marketing description based on name and ingredients
// @Tags         ai
// @Accept       json
// @Produce      json
// @Param        input body map[string]interface{} true "JSON input: {name: string, ingredients: []string}"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /menu/ai/generate-description [post]
func (c *MenuController) GenerateDescription(ctx *gin.Context) {
	var input struct {
		Name        string   `json:"name"`
		Ingredients []string `json:"ingredients"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	desc, err := c.service.GenerateDescriptionAI(input.Name, input.Ingredients)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "AI Service Error: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"generated_description": desc,
	})
}

// GetRecommendation godoc
// @Summary      AI Menu Recommendation
// @Description  Get menu recommendations based on user preference using Gemini AI
// @Tags         menu
// @Accept       json
// @Produce      json
// @Param        request body model.RecommendationRequest true "User Preference"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /menu/recommend [post]
func (c *MenuController) GetRecommendation(ctx *gin.Context) {
	var req model.RecommendationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recommendation, err := c.service.GetRecommendationAI(req.Preference)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "AI Service unavailable: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"recommendation": recommendation,
		},
	})
}
