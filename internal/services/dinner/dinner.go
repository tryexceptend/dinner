// Основная логика приложения.
// Сервис отдает массив из блюд для ужина.
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

// Доступ к списку доступных блюд
type FoodProvider interface {
	GetFoods() ([]models.Food, error)
}

// Доступ к истории запросов пользователей
type HistoryProvider interface {
	SaveRequest(userId int64) error
	IsLimit(userId int64) (bool, error)
}

// New - конструктор сервиса
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

// GetRandomDinner отдает массив блюд на ужин для юзера userId.
func (d *Dinner) GetRandomDinner(userId int64) ([]models.Food, error) {
	const op = "Dinner.GetRandomDinner"

	log := d.log.With(
		slog.String("op", op),
	)
	// проверка на лимит запросов
	limit, err := d.historyProvider.IsLimit(userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !limit {
		return nil, fmt.Errorf("%s: %w", op, services.ErrAttemptLimitExceeded)
	}

	// Запрос списка доступных блюд
	foods, err := d.foodProvider.GetFoods()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//Сохранение запроса пользователя в истории
	err = d.historyProvider.SaveRequest(userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("select dinner save reques")

	// Проверка, что список блюд не пустой
	if len(foods) == 0 {
		return nil, fmt.Errorf("%s: %w", op, services.ErrEmptyFood)
	}

	// Подучение случайного блюда
	rndPos := rand.IntN(len(foods))
	food := make([]models.Food, 1, 2)
	food[0] = foods[rndPos]

	// В зависимости от типа блюда отдаем 1 блюдо или ищем гранир к мясу
	switch food[0].Category {
	case models.Soup:
		return food, nil
	case models.Salad:
		return food, nil
	case models.Meat:
		sideDishes := GetSideDishes(&foods)
		if len(sideDishes) == 0 {
			return food, nil
		}
		food = append(food, sideDishes[rand.IntN(len(sideDishes))])
		return food, nil
	case models.SideDish:
		meats := GetMeats(&foods)
		if len(meats) == 0 {
			return food, nil
		}
		food = append(food, meats[rand.IntN(len(meats))])
		return food, nil
	}
	return nil, services.ErrEmptyFood
}

// GetSideDishes ищет гарнир к мясу
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

// GetMeats ищет мясо к гарниру
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
