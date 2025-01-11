// Пакет для телеграм бота
package app

import (
	"dinner/internal/config"
	dinnerservice "dinner/internal/services/dinner"
	storagesqlite "dinner/internal/storages/sqlite"
	telegrambot "dinner/internal/telegramBot"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

// Приложение бота
type App struct {
	Bot *telegrambot.TelegramBot
}

// New при помощи конфига создает и настраивает бота
func New(log *slog.Logger, token string, config *config.Config) *App {
	// Создает слой работы с БД в виде storage
	storage, err := storagesqlite.New(log, config.StoragePath)
	if err != nil {
		panic(err)
	}
	// Создает сервисный слой в виде сервиса dinner
	dinner := dinnerservice.New(log, storage, storage)
	// Создает инфраструктурный слой в вибе бота
	bot := telegrambot.New(log, token, config.Timeout, dinner)
	return &App{
		Bot: bot,
	}
}
