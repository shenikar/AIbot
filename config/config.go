package config

import (
	"os"
)

type Config struct {
	BotToken  string
	OpenAIkey string
	RedisAddr string
}

func LoadConfig() Config {
	return Config{
		BotToken:  os.Getenv("BOT_TOKEN"),
		OpenAIkey: os.Getenv("OPENAI_KEY"),
		RedisAddr: os.Getenv("REDIS_ADDR"),
	}
}
