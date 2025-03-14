package tgbot

import (
	"AIbot/internal/openai"
	"AIbot/internal/storage"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TGBot struct {
	api           *tgbotapi.BotAPI
	openAIService *openai.OpenAIService
	storage       *storage.RedisStorage
}

func NewTGBot(token string, ai *openai.OpenAIService, storage *storage.RedisStorage) (*TGBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Println("Ошибка создания Telegram-бота:", err)
		return nil, err
	}

	log.Println("Бот успешно создан! Запущен под именем:", bot.Self.UserName)
	return &TGBot{
		api:           bot,
		openAIService: ai,
		storage:       storage,
	}, nil
}

func (b *TGBot) StartPolling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		userID := update.Message.Chat.ID
		userMsg := update.Message.Text
		history := b.storage.GetHistory(userID)
		history = append(history, userMsg)
		aiResponse, err := b.openAIService.GetAIResponse(history)
		if err != nil {
			log.Printf("Ошибка при обработке запроса для пользователя %d: %v", userID, err)
			aiResponse = "Ошибка обработки запроса."
		}
		b.storage.SaveMessage(userID, userMsg)
		b.storage.SaveMessage(userID, aiResponse)

		msg := tgbotapi.NewMessage(userID, aiResponse)
		b.api.Send(msg)
	}
}
