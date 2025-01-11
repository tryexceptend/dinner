package dinner

import (
	"dinner/internal/domain/models"
	"dinner/internal/services"
	dinnerservice "dinner/internal/services/dinner"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFoodProvider struct {
	mock.Mock
}

func (m *MockFoodProvider) GetFoods() ([]models.Food, error) {
	args := m.Called()
	return args.Get(0).([]models.Food), args.Error(1)
}

type MockHistoryProvider struct {
	mock.Mock
}

func (m *MockHistoryProvider) SaveRequest(userId int64) error {
	args := m.Called(userId)
	return args.Error(0)
}
func (m *MockHistoryProvider) IsLimit(userId int64) (bool, error) {
	args := m.Called(userId)
	return args.Get(0).(bool), args.Error(1)
}
func TestGetSideDishesEmpty(t *testing.T) {
	foods := []models.Food{}
	actual := dinnerservice.GetSideDishes(&foods)
	if len(actual) != 0 {
		t.Errorf("GetSideDishes empty return not empty array")
	}
}

func TestGetRandomDinnerIsLimit(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	mockFoodProvider := new(MockFoodProvider)
	mockFoodProvider.On("GetFoods").Return(nil, nil)

	mockHistoryProvider := new(MockHistoryProvider)
	mockHistoryProvider.On("SaveRequest", mock.Anything).Return(nil)
	mockHistoryProvider.On("IsLimit", mock.Anything).Return(false, nil)

	dinnerService := dinnerservice.New(log, mockFoodProvider, mockHistoryProvider)
	_, err := dinnerService.GetRandomDinner(1)

	if !errors.Is(err, services.ErrAttemptLimitExceeded) {
		t.Errorf("return incorrect error: " + err.Error())
	}
}

func TestGetRandomDinnerNilFood(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	var foodNil []models.Food
	mockFoodProvider := new(MockFoodProvider)
	mockFoodProvider.On("GetFoods").Return(foodNil, nil)

	mockHistoryProvider := new(MockHistoryProvider)
	mockHistoryProvider.On("SaveRequest", mock.Anything).Return(nil)
	mockHistoryProvider.On("IsLimit", mock.Anything).Return(true, nil)

	dinnerService := dinnerservice.New(log, mockFoodProvider, mockHistoryProvider)
	_, err := dinnerService.GetRandomDinner(1)

	if !errors.Is(err, services.ErrEmptyFood) {
		t.Errorf("return incorrect error: " + err.Error())
	}
}
func TestGetRandomDinnerEmptyFood(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	mockFoodProvider := new(MockFoodProvider)
	mockFoodProvider.On("GetFoods").Return([]models.Food{}, nil)

	mockHistoryProvider := new(MockHistoryProvider)
	mockHistoryProvider.On("SaveRequest", mock.Anything).Return(nil)
	mockHistoryProvider.On("IsLimit", mock.Anything).Return(true, nil)

	dinnerService := dinnerservice.New(log, mockFoodProvider, mockHistoryProvider)
	_, err := dinnerService.GetRandomDinner(1)

	if !errors.Is(err, services.ErrEmptyFood) {
		t.Errorf("return incorrect error: " + err.Error())
	}
}

func TestGetRandomDinnerCorrectFood(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	tests := []struct {
		name     string
		foods    []models.Food
		expected int
	}{
		{
			name: "one food salad",
			foods: []models.Food{
				models.Food{
					Name:     "Salad1",
					Category: models.Salad,
				},
			},
			expected: 1,
		},
		{
			name: "two food salad",
			foods: []models.Food{
				models.Food{
					Name:     "Salad1",
					Category: models.Salad,
				},
				models.Food{
					Name:     "Salad2",
					Category: models.Salad,
				},
			},
			expected: 1,
		},
		{
			name: "one food soup",
			foods: []models.Food{
				models.Food{
					Name:     "Soup1",
					Category: models.Soup,
				},
			},
			expected: 1,
		},
		{
			name: "two food soup",
			foods: []models.Food{
				models.Food{
					Name:     "Soup1",
					Category: models.Soup,
				},
				models.Food{
					Name:     "Soup2",
					Category: models.Soup,
				},
			},
			expected: 1,
		},
		{
			name: "soup and salat",
			foods: []models.Food{
				models.Food{
					Name:     "Soup1",
					Category: models.Soup,
				},
				models.Food{
					Name:     "Salad1",
					Category: models.Salad,
				},
			},
			expected: 1,
		},
		{
			name: "two foods",
			foods: []models.Food{
				models.Food{
					Name:     "Meat1",
					Category: models.Meat,
				},
				models.Food{
					Name:     "SideDish1",
					Category: models.SideDish,
				},
			},
			expected: 2,
		},
		{
			name: "two meats and one sidedish",
			foods: []models.Food{
				models.Food{
					Name:     "Meat1",
					Category: models.Meat,
				},
				models.Food{
					Name:     "Meat2",
					Category: models.Meat,
				},
				models.Food{
					Name:     "SideDish1",
					Category: models.SideDish,
				},
			},
			expected: 2,
		},
		{
			name: "one meat and two sidedishes",
			foods: []models.Food{
				models.Food{
					Name:     "Meat1",
					Category: models.Meat,
				},
				models.Food{
					Name:     "SideDish1",
					Category: models.SideDish,
				},
				models.Food{
					Name:     "SideDish2",
					Category: models.SideDish,
				},
			},
			expected: 2,
		},
		{
			name: "two meats",
			foods: []models.Food{
				models.Food{
					Name:     "Meat1",
					Category: models.Meat,
				},
				models.Food{
					Name:     "Meat2",
					Category: models.Meat,
				},
			},
			expected: 1,
		},
		{
			name: "two sideDishes",
			foods: []models.Food{
				models.Food{
					Name:     "SideDish1",
					Category: models.SideDish,
				},
				models.Food{
					Name:     "SideDish2",
					Category: models.SideDish,
				},
			},
			expected: 1,
		},
	}

	mockHistoryProvider := new(MockHistoryProvider)
	mockHistoryProvider.On("SaveRequest", mock.Anything).Return(nil)
	mockHistoryProvider.On("IsLimit", mock.Anything).Return(true, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFoodProvider := new(MockFoodProvider)
			mockFoodProvider.On("GetFoods").Return(tt.foods, nil)
			dinnerService := dinnerservice.New(log, mockFoodProvider, mockHistoryProvider)

			foods, err := dinnerService.GetRandomDinner(1)
			assert.Nil(t, err)
			assert.Len(t, foods, tt.expected)
		})
	}
}

func TestGetRandomDinnerCorrectTwoFood(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	tests := []struct {
		name     string
		foods    []models.Food
		expected int
	}{
		{
			name: "two foods",
			foods: []models.Food{
				models.Food{
					Name:     "Meat1",
					Category: models.Meat,
				},
				models.Food{
					Name:     "SideDish1",
					Category: models.SideDish,
				},
			},
			expected: 2,
		},
		{
			name: "two meats and one sidedish",
			foods: []models.Food{
				models.Food{
					Name:     "Meat1",
					Category: models.Meat,
				},
				models.Food{
					Name:     "Meat2",
					Category: models.Meat,
				},
				models.Food{
					Name:     "SideDish1",
					Category: models.SideDish,
				},
			},
			expected: 2,
		},
		{
			name: "one meat and two sidedishes",
			foods: []models.Food{
				models.Food{
					Name:     "Meat1",
					Category: models.Meat,
				},
				models.Food{
					Name:     "SideDish1",
					Category: models.SideDish,
				},
				models.Food{
					Name:     "SideDish2",
					Category: models.SideDish,
				},
			},
			expected: 2,
		},
		{
			name: "two meats and two sidedishes",
			foods: []models.Food{
				models.Food{
					Name:     "Meat1",
					Category: models.Meat,
				},
				models.Food{
					Name:     "Meat2",
					Category: models.Meat,
				},
				models.Food{
					Name:     "SideDish1",
					Category: models.SideDish,
				},
				models.Food{
					Name:     "SideDish2",
					Category: models.SideDish,
				},
			},
			expected: 2,
		},
	}

	mockHistoryProvider := new(MockHistoryProvider)
	mockHistoryProvider.On("SaveRequest", mock.Anything).Return(nil)
	mockHistoryProvider.On("IsLimit", mock.Anything).Return(true, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFoodProvider := new(MockFoodProvider)
			mockFoodProvider.On("GetFoods").Return(tt.foods, nil)
			dinnerService := dinnerservice.New(log, mockFoodProvider, mockHistoryProvider)

			foods, err := dinnerService.GetRandomDinner(1)
			assert.Nil(t, err)
			assert.Len(t, foods, tt.expected)
			assert.True(t, (foods[0].Category == models.Meat && foods[1].Category == models.SideDish) || (foods[1].Category == models.Meat && foods[0].Category == models.SideDish))
		})
	}
}
