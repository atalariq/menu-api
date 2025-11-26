package main

import (
	"fmt"
	"log"
	"os"

	"atalariq/menu-api/internal/controller"
	"atalariq/menu-api/internal/model"
	"atalariq/menu-api/internal/repository"
	"atalariq/menu-api/internal/service"

	_ "atalariq/menu-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title           Menu API
// @version         1.0
// @description     API for Menu Management
// @host            localhost:8080
// @BasePath        /
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	// 1. DB Connection
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "menu.db" // Fallback untuk local development
	}
	db, _ := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

	if err := db.AutoMigrate(&model.Menu{}); err != nil {
		log.Fatal("Failed to migrate database scheme:", err)
	}

	// 2. Dependency Injection
	menuRepository := repository.NewMenuRepository(db)
	menuService := service.NewMenuService(menuRepository)
	menuController := controller.NewMenuController(menuService)

	// 3. Router
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/menu")
	{
		api.POST("", menuController.Create)
		api.GET("", menuController.FindAll)
		api.GET("/:id", menuController.GetByID)
		api.PUT("/:id", menuController.Update)
		api.DELETE("/:id", menuController.Delete)
		api.GET("/group-by-category", menuController.GroupByCategory)
		api.GET("/search", menuController.FindAll)

		// AI Routes
		api.POST("/generate-description", menuController.GenerateDescription)
		api.POST("/recommend", menuController.GetRecommendation)
	}

	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
