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

		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
			continue
		}

		history := b.storage.GetHistory(userID)
		history = append(history, userMsg)
		aiResponse, err := b.openAIService.GetAIResponse(history)
		if err != nil {
			log.Printf("Ошибка при обработке запроса для пользователя %d: %v", userID, err)
			aiResponse = "Ошибка обработки запроса."
		}
		b.storage.SaveMessage(userID, userMsg)
		b.storage.SaveMessage(userID, aiResponse)

		b.sendMessage(userID, aiResponse)
	}
}

func (b *TGBot) handleCommand(msg *tgbotapi.Message) {
	log.Printf("Получена команда: %s", msg.Command())

	switch msg.Command() {
	case "image":
		b.handleImageCommand(msg)
	default:
		b.sendMessage(msg.Chat.ID, "Неизвестная команда.")
	}
}

func (b *TGBot) handleImageCommand(msg *tgbotapi.Message) {
	prompt := msg.CommandArguments()

	if prompt == "" {
		b.sendMessage(msg.Chat.ID, "Пожалуйста, укажите описание для генерации изображения.")
		return
	}

	imageURLs, err := b.openAIService.GenerateImage(prompt)
	if err != nil {
		b.sendMessage(msg.Chat.ID, "Ошибка при генерации изображения: "+err.Error())
		return
	}

	if len(imageURLs) > 0 {
		photo := tgbotapi.NewPhoto(msg.Chat.ID, tgbotapi.FileURL(imageURLs[0]))
		b.api.Send(photo)
	} else {
		b.sendMessage(msg.Chat.ID, "Не удалось сгенерировать изображение.")
	}
}

func (b *TGBot) sendMessage(chatID int64, text string) {
	log.Println("Отправка сообщения пользователю", chatID, ":", text)
	msg := tgbotapi.NewMessage(chatID, text)
	b.api.Send(msg)
}
