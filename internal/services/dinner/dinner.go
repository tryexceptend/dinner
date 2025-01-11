package dinnerservice

import (
	"dinner/internal/domain/models"
	"dinner/internal/services"
	"fmt"
	"log/slog"
	"math/rand/v2"
)

type Dinner struct {
	log             *slog.Logger
	foodProvider    FoodProvider
	historyProvider HistoryProvider
}

type FoodProvider interface {
	GetFoods() ([]models.Food, error)
}

type HistoryProvider interface {
	SaveRequest(userId int64) error
	IsLimit(userId int64) (bool, error)
}

func New(
	log *slog.Logger,
	foodProvider FoodProvider,
	historyProvider HistoryProvider,
) *Dinner {
	return &Dinner{
		log:             log,
		foodProvider:    foodProvider,
		historyProvider: historyProvider,
	}
}

// IsAdmin checks if user is admin.
func (d *Dinner) GetRandomDinner(userId int64) ([]models.Food, error) {
	const op = "Dinner.GetRandomDinner"

	log := d.log.With(
		slog.String("op", op),
	)

	limit, err := d.historyProvider.IsLimit(userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !limit {
		return nil, fmt.Errorf("%s: %w", op, services.ErrAttemptLimitExceeded)
	}

	foods, err := d.foodProvider.GetFoods()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = d.historyProvider.SaveRequest(userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(foods) == 0 {
		return nil, fmt.Errorf("%s: %w", op, services.ErrEmptyFood)
	}

	rndPos := rand.IntN(len(foods))

	food := make([]models.Food, 1, 2)
	food[0] = foods[rndPos]

	switch food[0].Category {
	case models.Soup:
		return food, nil
	case models.Salad:
		return food, nil
	}
	if food[0].Category == models.Meat {
		sideDishes := GetSideDishes(&foods)
		if len(sideDishes) == 0 {
			return food, nil
		}
		food = append(food, sideDishes[rand.IntN(len(sideDishes))])
		return food, nil
	}
	if food[0].Category == models.SideDish {
		meats := GetMeats(&foods)
		if len(meats) == 0 {
			return food, nil
		}
		food = append(food, meats[rand.IntN(len(meats))])
		return food, nil
	}

	log.Info("select dinner save reques", slog.Any("food", food))
	return food, nil
}

func GetSideDishes(foods *[]models.Food) []models.Food {
	res := make([]models.Food, 0)
	if len(*foods) == 0 {
		return res
	}
	for _, value := range *foods {
		if value.Category == models.SideDish {
			res = append(res, value)
		}
	}
	return res
}

func GetMeats(foods *[]models.Food) []models.Food {
	res := make([]models.Food, 0)
	if len(*foods) == 0 {
		return res
	}
	for _, value := range *foods {
		if value.Category == models.Meat {
			res = append(res, value)
		}
	}
	return res
}
