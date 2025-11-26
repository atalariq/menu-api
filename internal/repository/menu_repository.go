// Package repository
package repository

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"atalariq/menu-api/internal/model"

	"gorm.io/gorm"
)

type MenuRepository interface {
	Create(menu *model.Menu) error
	FindAll(filter model.MenuFilter) ([]model.Menu, model.PaginationResponse, error)
	FindByID(id uint) (model.Menu, error)
	Update(menu *model.Menu) error
	Delete(id uint) error
	GroupBy(mode string, limit int) (interface{}, error)
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db}
}

func (r *menuRepository) Create(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

func (r *menuRepository) FindAll(filter model.MenuFilter) ([]model.Menu, model.PaginationResponse, error) {
	var menus []model.Menu
	var total int64

	db := r.db.Model(&model.Menu{})

	// Apply Filters
	if filter.Query != "" {
		db = db.Where("name LIKE ? OR description LIKE ?", "%"+filter.Query+"%", "%"+filter.Query+"%")
	}
	if filter.Category != "" {
		db = db.Where("category = ?", filter.Category)
	}
	if filter.MinPrice > 0 {
		db = db.Where("price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		db = db.Where("price <= ?", filter.MaxPrice)
	}
	if filter.MaxCal > 0 {
		db = db.Where("calories <= ?", filter.MaxCal)
	}

	// Count Total (for pagination)
	db.Count(&total)

	// Sorting
	if filter.Sort != "" {
		parts := strings.Split(filter.Sort, ":")
		if len(parts) == 2 {
			db = db.Order(fmt.Sprintf("%s %s", parts[0], parts[1]))
		}
	} else {
		db = db.Order("created_at desc")
	}

	// Pagination Logic
	offset := (filter.Page - 1) * filter.PerPage
	db = db.Limit(filter.PerPage).Offset(offset)

	err := db.Find(&menus).Error

	totalPages := int(math.Ceil(float64(total) / float64(filter.PerPage)))
	pagination := model.PaginationResponse{
		Total:      total,
		Page:       filter.Page,
		PerPage:    filter.PerPage,
		TotalPages: totalPages,
	}

	return menus, pagination, err
}

func (r *menuRepository) FindByID(id uint) (model.Menu, error) {
	var menu model.Menu
	err := r.db.First(&menu, id).Error
	return menu, err
}

func (r *menuRepository) Update(menu *model.Menu) error {
	return r.db.Save(menu).Error
}

func (r *menuRepository) Delete(id uint) error {
	return r.db.Delete(&model.Menu{}, id).Error
}

func (r *menuRepository) GroupBy(mode string, limit int) (interface{}, error) {
	if mode == "count" {
		type Result struct {
			Category string
			Count    int
		}
		var results []Result

		err := r.db.Model(&model.Menu{}).
			Select("category, count(*) as count").
			Group("category").
			Scan(&results).Error
		if err != nil {
			return nil, err
		}

		output := make(map[string]int)
		for _, res := range results {
			output[res.Category] = res.Count
		}
		return output, nil
	}

	if mode == "list" {
		var menus []model.Menu

		if err := r.db.Order("category asc, name asc").Find(&menus).Error; err != nil {
			return nil, err
		}

		// Key: Category name, Value: Slice of Menu
		grouped := make(map[string][]model.Menu)

		for _, menu := range menus {
			if limit > 0 {
				if len(grouped[menu.Category]) >= limit {
					continue
				}
			}

			grouped[menu.Category] = append(grouped[menu.Category], menu)
		}

		return grouped, nil
	}

	return nil, errors.New("invalid mode")
}
