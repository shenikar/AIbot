package tgbot

import (
	"AIbot/internal/openai"
	"AIbot/internal/storage"
	"log"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TGBot struct {
	api                   *tgbotapi.BotAPI
	openAIService         *openai.OpenAIService
	storage               *storage.RedisStorage
	waitingForImagePrompt map[int64]bool // Флаг ожидания описания изображения
	mu                    sync.Mutex
}

func NewTGBot(token string, ai *openai.OpenAIService, storage *storage.RedisStorage) (*TGBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Println("Ошибка создания Telegram-бота:", err)
		return nil, err
	}

	log.Println("Бот успешно создан! Запущен под именем:", bot.Self.UserName)
	return &TGBot{
		api:                   bot,
		openAIService:         ai,
		storage:               storage,
		waitingForImagePrompt: make(map[int64]bool),
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

		// Проверяем, ожидает ли бот описание для генерации изображения
		b.mu.Lock()
		waiting := b.waitingForImagePrompt[userID]
		b.mu.Unlock()

		if waiting {
			b.handleImageGeneration(userID, userMsg)
			continue
		}

		// Если сообщение - это команда, обрабатываем команду
		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
			continue
		}

		// Если текст начинается с "нарисуй", обрабатываем как генерацию изображения
		if strings.HasPrefix(strings.ToLower(userMsg), "нарисуй") {
			b.handleImageGeneration(userID, strings.TrimSpace(strings.TrimPrefix(userMsg, "нарисуй")))
			continue
		}

		// Обычные текстовые сообщения
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

// Обработка команды /image
func (b *TGBot) handleImageRequest(msg *tgbotapi.Message) {
	b.mu.Lock()
	b.waitingForImagePrompt[msg.Chat.ID] = true
	b.mu.Unlock()

	b.sendMessage(msg.Chat.ID, "Пожалуйста, отправьте описание для генерации изображения. Например: 'кот на крыше'.")
}

// Генерация изображения
func (b *TGBot) handleImageGeneration(userID int64, prompt string) {
	// Сбрасываем флаг ожидания описания
	b.mu.Lock()
	b.waitingForImagePrompt[userID] = false
	b.mu.Unlock()

	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		b.sendMessage(userID, "Пожалуйста, укажите описание для генерации изображения. Например: 'кот на крыше'.")
		return
	}

	b.sendMessage(userID, "Генерирую изображение, пожалуйста, подождите...")

	imageURLs, err := b.openAIService.GenerateImage(prompt)
	if err != nil {
		b.sendMessage(userID, "Ошибка при генерации изображения: "+err.Error())
		return
	}

	if len(imageURLs) > 0 {
		photo := tgbotapi.NewPhoto(userID, tgbotapi.FileURL(imageURLs[0]))
		b.api.Send(photo)
	} else {
		b.sendMessage(userID, "Не удалось сгенерировать изображение.")
	}
}

func (b *TGBot) handleCommand(msg *tgbotapi.Message) {
	log.Printf("Получена команда: %s", msg.Command())

	switch msg.Command() {
	case "start":
		b.handleStart(msg)
	case "image":
		b.handleImageRequest(msg)
	default:
		b.sendMessage(msg.Chat.ID, "Неизвестная команда.")
	}
}

// Обработка команды /start
func (b *TGBot) handleStart(msg *tgbotapi.Message) {
	welcomeText := "Привет! Я бот, который может помочь с генерацией изображений и ответами на ваши текстовые сообщения.\n\n" +
		"1. Для генерации изображения используйте команду /image и отправьте описание.\n" +
		"   Например: 'кот на крыше'.\n" +
		"2. Можно сразу написать 'нарисуй' + описание.\n" +
		"   Например: 'нарисуй дракона в небе'.\n" +
		"3. Я также могу отвечать на текстовые сообщения, просто напишите мне!"
	b.sendMessage(msg.Chat.ID, welcomeText)
}

func (b *TGBot) sendMessage(chatID int64, text string) {
	log.Println("Отправка сообщения пользователю", chatID, ":", text)
	msg := tgbotapi.NewMessage(chatID, text)
	b.api.Send(msg)
}
