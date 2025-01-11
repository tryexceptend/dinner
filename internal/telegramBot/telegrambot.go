package telegrambot

import (
	"dinner/internal/services"
	dinnerservice "dinner/internal/services/dinner"
	"errors"
	"log/slog"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	log     *slog.Logger
	token   string
	timeout int
	dinner  *dinnerservice.Dinner
}

// New Конструктор бота
// log *slog.Logger - логгер
// token string - токен бота (берется из переменной окружения)
// timeout int - таймаут
// dinner *dinnerservice.Dinner - сервис, который генерит что приготовить на ужин
func New(log *slog.Logger, token string, timeout int, dinner *dinnerservice.Dinner) *TelegramBot {
	// TODO: проверка валидности токена
	return &TelegramBot{
		log:     log,
		token:   token,
		timeout: timeout,
		dinner:  dinner,
	}
}

// Run основной поток бота
func (b *TelegramBot) Run() {
	const op = "TelegramBot.Run"
	log := b.log.With(slog.String("op", op))
	bot, err := tgbotapi.NewBotAPI(b.token)
	if err != nil {
		panic(err)
	}

	bot.Debug = false

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = b.timeout
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {

		if update.Message == nil {
			continue
		}
		// Обработка команды /dinner
		if update.Message.Text == "/dinner" {
			err := b.DinnerCommand(bot, update.Message)
			if err != nil {
				continue
			}
			log.Info("apply command '/dinner'")
		}
	}
}

// DinnerCommand запрашивет у сервиса блюда на ужин.
func (b *TelegramBot) DinnerCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {
	if message.Text != "/dinner" {
		return nil
	}
	const op = "TelegramBot.DinnerCommand"
	log := b.log.With(slog.String("op", op))
	// Получение блюд
	foods, err := b.dinner.GetRandomDinner(message.From.ID)
	if err != nil {
		// Превышен лимит запросов
		if errors.Is(err, services.ErrAttemptLimitExceeded) {
			b.log.Debug("user attempt limit exceeded", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			msg := tgbotapi.NewMessage(message.Chat.ID, "Лимит попыток исчерпан")
			if _, err := bot.Send(msg); err != nil {
				log.Error("send message error", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
			}
		}

		log.Error("get random dinner error", slog.Any("error", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())}))
		return err
	}
	// Нет блюд
	if len(foods) == 0 {
		log.Error("get random dinner error", slog.Any("error", slog.Attr{Key: "error", Value: slog.StringValue(services.ErrEmptyFood.Error())}))
		return services.ErrEmptyFood
	}
	// Формирование ответного сообщения
	msgFood := foods[0].Name
	for i := 1; i < len(foods); i++ {
		msgFood = msgFood + " и " + strings.ToLower(foods[i].Name)
	}
	// Отправка сообщения пользователю
	msg := tgbotapi.NewMessage(message.Chat.ID, msgFood)
	if _, err := bot.Send(msg); err != nil {
		log.Error("send message error", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})
	}
	return nil
}
