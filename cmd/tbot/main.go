// Точка запуска бота
package main

import (
	"dinner/internal/app"
	"dinner/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Получаем конфигурацию
	cfg := config.MustLoadConfig()

	// Создаем логгер
	logger := mustSetupLogger(cfg.Env)
	logger.Debug("Start application", slog.Any("cfg", cfg))

	// Инициализируем приложение
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		panic("bot tokent is empty")
	}
	application := app.New(logger, token, cfg)
	// Запускаем бота на прослушивание в горутине
	go application.Bot.Run()

	// Ожидаем от системы команды на остановку
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	logger.Warn("application stopped", slog.String("signal", sign.String()))
}

// mustSetupLogger настраивает логгер
func mustSetupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case config.EnvLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.EnvDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.EnvProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		panic("logger not configured")
	}
	return logger
}
