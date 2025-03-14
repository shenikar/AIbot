package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr string) *RedisStorage {
	client := redis.NewClient(&redis.Options{Addr: addr})

	// Проверяем соединение с Redis
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}

	log.Println("Redis успешно подключен!")

	return &RedisStorage{client: client}
}

func (r *RedisStorage) SaveMessage(userID int64, msg string) {
	ctx := context.Background()
	key := fmt.Sprintf("chat:%d", userID)

	err := r.client.RPush(ctx, key, msg).Err()
	if err != nil {
		log.Printf("Ошибка сохранения сообщения в Redis: %v", err)
		return
	}

	r.client.LTrim(ctx, key, -10, -1) // Оставляем только последние 10 сообщений
}

func (r *RedisStorage) GetHistory(userID int64) []string {
	ctx := context.Background()
	key := fmt.Sprintf("chat:%d", userID)

	messages, err := r.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		log.Printf("Ошибка получения истории сообщений из Redis: %v", err)
		return nil
	}

	return messages
}
