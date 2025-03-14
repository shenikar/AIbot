package main

import (
	"AIbot/config"
	"AIbot/internal/openai"
	"AIbot/internal/storage"
	"AIbot/internal/tgbot"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	storage := storage.NewRedisStorage(cfg.RedisAddr)

	// Передаем ключ ProxyAPI и URL ProxyAPI
	aiService := openai.NewOpenAIService(cfg.ProxyAPIKey, cfg.ProxyAPIURL)

	bot, err := tgbot.NewTGBot(cfg.BotToken, aiService, storage)
	if err != nil {
		log.Fatal("Ошибка запуска бота:", err)
	}

	bot.StartPolling()
}
