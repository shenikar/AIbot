package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken    string
	ProxyAPIKey string
	ProxyAPIURL string
	RedisAddr   string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка загрузки .env файла")
	}
	cfg := Config{
		BotToken:    os.Getenv("BOT_TOKEN"),
		ProxyAPIKey: os.Getenv("ProxyAPIKey"),
		ProxyAPIURL: os.Getenv("ProxyAPIURL"),
		RedisAddr:   os.Getenv("REDIS_ADDR"),
	}
	return cfg
}
