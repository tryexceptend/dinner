package services

import "errors"

// Ошибки на доменном уровне
var (
	// Превышен лимит запросов от юзера
	ErrAttemptLimitExceeded = errors.New("user attempt limit exceeded")
	// Не удалось сформировать ужин
	ErrEmptyFood = errors.New("food is empty")
)
