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
//
// @Summary    Create a new menu
// @Description  Create a new menu item with ingredients
// @Tags     menu
// @Accept     json
// @Produce    json
// @Param      menu  body    model.Menu          true  "Menu Request"
// @Success    201   {object}  model.MenuSuccessResponse "Typed Response"
// @Failure    400  {object}  model.ErrorResponse  "Validation Error"
// @Failure    500  {object}  model.ErrorResponse  "Server Error"
// @Router     /menu [post]
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

// GetList godoc
// @Summary      List menus (Browsing)
// @Description  Get menu list with filtering, sorting, and pagination
// @Tags         menu
// @Produce      json
// @Param        category   query     string  false  "Filter by category"
// @Param        min_price  query     number  false  "Minimum price"
// @Param        max_price  query     number  false  "Maximum price"
// @Param        max_cal    query     int     false  "Maximum calories"
// @Param        sort       query     string  false  "Sort (e.g., price:asc)"
// @Param        page       query     int     false  "Page number (default 1)"
// @Param        per_page   query     int     false  "Items per page (default 10)"
// @Success      200        {object}  model.MenuPaginationResponse
// @Failure      400        {object}  model.ErrorResponse
// @Router       /menu [get]
func (c *MenuController) GetList(ctx *gin.Context) {
	var params model.MenuQueryRequest

	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := model.MenuFilter{
		Category: params.Category,
		MinPrice: params.MinPrice,
		MaxPrice: params.MaxPrice,
		MaxCal:   params.MaxCal,
		Sort:     params.Sort,
		Page:     params.Page,
		PerPage:  params.PerPage,
	}

	result, err := c.service.GetList(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// Search godoc
// @Summary      Search menus
// @Description  Search menu by name or description (Full Text Search intent)
// @Tags         menu
// @Produce      json
// @Param        q          query     string  false   "Search keyword"
// @Param        category   query     string  false  "Filter by category"
// @Param        min_price  query     number  false  "Minimum price"
// @Param        max_price  query     number  false  "Maximum price"
// @Param        sort       query     string  false  "Sort (e.g., price:asc)"
// @Param        page       query     int     false  "Page number (default 1)"
// @Param        per_page   query     int     false  "Items per page (default 10)"
// @Success      200        {object}  model.MenuPaginationResponse
// @Failure      400        {object}  model.ErrorResponse
// @Router       /menu/search [get]
func (c *MenuController) Search(ctx *gin.Context) {
	var params model.MenuQueryRequest

	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Make param q not mandatory
	// if params.Q == "" {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Query param 'q' is required for search"})
	// 	return
	// }

	filter := model.MenuFilter{
		Query:    params.Q,
		Category: params.Category,
		MinPrice: params.MinPrice,
		MaxPrice: params.MaxPrice,
		Sort:     params.Sort,
		Page:     params.Page,
		PerPage:  params.PerPage,
	}

	result, err := c.service.GetList(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetByID godoc
//
// @Summary    Get menu detail
// @Description  Get details of a specific menu item by ID
// @Tags     menu
// @Produce    json
// @Param      id  path    int             true  "Menu ID"
// @Success    200 {object}  model.MenuDetailResponse  "Typed Response"
// @Failure    400 {object}  model.ErrorResponse  "Invalid ID"
// @Failure    404 {object}  model.ErrorResponse  "Menu Not Found"
// @Router     /menu/{id} [get]
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
//
// @Summary    Update menu
// @Description  Update an existing menu item
// @Tags     menu
// @Accept     json
// @Produce    json
// @Param      id    path      int             true  "Menu ID"
// @Param      menu  body      model.Menu          true  "Update Data"
// @Success    200   {object}  model.MenuSuccessResponse "Typed Response"
// @Failure    400   {object}  model.ErrorResponse  "Invalid ID"
// @Failure    404   {object}  model.ErrorResponse  "Menu Not Found"
// @Router     /menu/{id} [put]
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
//
// @Summary    Delete menu
// @Description  Delete a menu item by ID
// @Tags     menu
// @Produce    json
// @Param      id  path    int true  "Menu ID"
// @Success    200 {object}  model.GeneralResponse
// @Failure    400   {object}  model.ErrorResponse  "Invalid ID"
// @Failure    404   {object}  model.ErrorResponse  "Menu Not Found"
// @Router     /menu/{id} [delete]
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
//
// @Summary    Group menus by category
// @Description  Get menu counts or lists grouped by category
// @Tags     menu
// @Produce    json
// @Param      mode      query   string  true  "Mode: 'count' or 'list'"
// @Param      per_category  query   int   false "Limit item per category (default 5)"
// @Success    200       {object}  map[string]any
// @Failure    400  {object}  model.ErrorResponse  "Invalid mode"
// @Failure    500  {object}  model.ErrorResponse  "Server Error"
// @Router     /menu/group-by-category [get]
func (c *MenuController) GroupByCategory(ctx *gin.Context) {
	limitPerCategory, err := strconv.Atoi(ctx.Query("per_category"))
	if err != nil {
		limitPerCategory = 5
	}

	mode := ctx.Query("mode")
	if mode != "count" && mode != "list" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mode. Use 'count' or 'list'"})
		return
	}

	result, err := c.service.GetGrouped(mode, limitPerCategory)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// GenerateDescriptionAI godoc
//
// @Summary    Generate Menu Description
// @Description  Use Gemini AI to create a marketing description based on name and ingredients
// @Tags       AI
// @Accept     json
// @Produce    json
// @Param      input body      model.GenerateDescriptionRequest  true  "Input Data"
// @Success    200   {object}  model.GenerateDescriptionResponse
// @Failure    400   {object}  model.ErrorResponse  "Invalid input format"
// @Failure    500   {object}  model.ErrorResponse  "AI service error"
// @Router     /menu/generate-description [post]
func (c *MenuController) GenerateDescription(ctx *gin.Context) {
	var input model.GenerateDescriptionRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	desc, err := c.service.GenerateDescription(input.Name, input.Ingredients)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "AI Service Error: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"generated_description": desc,
	})
}

// GetRecommendations godoc
//
// @Summary      Get Menu Recommendations
// @Description  Get menu recommendations based on user preference using Gemini AI
// @Tags       AI
// @Accept     json
// @Produce    json
// @Param      request body    model.RecommendationRequest     true  "User Preference"
// @Success    200   {object}  model.RecommendationListResponse  "Typed Response"
// @Failure    400  {object}  model.ErrorResponse
// @Failure    502  {object}  model.ErrorResponse  "AI service unavailable"
// @Router     /menu/recommendations [post]
func (c *MenuController) GetRecommendations(ctx *gin.Context) {
	var request model.RecommendationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recommendations, err := c.service.GetRecommendations(request)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "AI Service unavailable: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"recommendations": recommendations,
	})
}
