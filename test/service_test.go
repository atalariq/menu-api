package test

import (
	"testing"

	"atalariq/menu-api/internal/model"
	"atalariq/menu-api/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository untuk memalsukan database
type MockRepository struct {
	mock.Mock
}

// Kita implementasikan interface repository.MenuRepository pada Mock ini
func (m *MockRepository) Create(menu *model.Menu) error {
	args := m.Called(menu)
	return args.Error(0)
}

func (m *MockRepository) FindAll(filter model.MenuFilter) ([]model.Menu, model.PaginationResponse, error) {
	// ... implementasi mock lainnya jika perlu
	return nil, model.PaginationResponse{}, nil
}

// ... (implementasikan method interface lainnya secara kosong/return args)
func (m *MockRepository) FindByID(id uint) (model.Menu, error)                { return model.Menu{}, nil }
func (m *MockRepository) Update(menu *model.Menu) error                       { return nil }
func (m *MockRepository) Delete(id uint) error                                { return nil }
func (m *MockRepository) GroupBy(mode string, limit int) (interface{}, error) { return nil, nil }

func TestCreateMenu_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	svc := service.NewMenuService(mockRepo)

	input := model.Menu{Name: "Burger", Price: 50000}

	// Expectation: Repo.Create dipanggil 1x, return nil error
	mockRepo.On("Create", &input).Return(nil)

	// Action
	result, err := svc.Create(input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Burger", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestCreateMenu_NegativePrice(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := service.NewMenuService(mockRepo)

	input := model.Menu{Name: "Burger", Price: -100}

	// Action (Logic error harus ditangkap sebelum masuk repo)
	_, err := svc.Create(input)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "price cannot be negative", err.Error())
	mockRepo.AssertNotCalled(t, "Create")
}
