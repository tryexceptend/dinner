package app

import (
	"dinner/internal/config"
	dinnerservice "dinner/internal/services/dinner"
	storagesqlite "dinner/internal/storages/sqlite"
	telegrambot "dinner/internal/telegramBot"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	Bot *telegrambot.TelegramBot
}

func New(log *slog.Logger, token string, config *config.Config) *App {
	storage, err := storagesqlite.New(log, config.StoragePath)
	if err != nil {
		panic(err)
	}

	dinner := dinnerservice.New(log, storage, storage)

	bot := telegrambot.New(log, token, config.Timeout, dinner)
	return &App{
		Bot: bot,
	}

}
