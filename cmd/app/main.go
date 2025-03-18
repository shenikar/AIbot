package main

import (
	"AIbot/config"
	"AIbot/internal/openai"
	"AIbot/internal/storage"
	"AIbot/internal/tgbot"
	"AIbot/pkg/logger"
	"log"
)

func main() {
	// Инициализируем логгер
	logg, err := logger.NewLogger("bot.log")
	if err != nil {
		log.Fatal("Ошибка инициализации логгера:", err)
	}

	logg.Info("Загрузка конфигурации...")
	cfg := config.LoadConfig()

	logg.Info("Подключение к Redis...")
	storage := storage.NewRedisStorage(cfg.RedisAddr)

	logg.Info("Создание сервиса OpenAI...")
	aiService := openai.NewOpenAIService(cfg.ProxyAPIKey, cfg.ProxyAPIURL)

	logg.Info("Запуск Telegram-бота...")
	bot, err := tgbot.NewTGBot(cfg.BotToken, aiService, storage)
	if err != nil {
		logg.Error("Ошибка запуска бота: " + err.Error())
		return
	}

	logg.Info("Бот успешно запущен, начинаю обработку сообщений...")
	bot.StartPolling()
}
