package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr string) *RedisStorage {
	return &RedisStorage{
		client: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func (r *RedisStorage) SaveMessage(userID int64, msg string) {
	ctx := context.Background()
	key := fmt.Sprintf("chat:%d", userID)
	r.client.RPush(ctx, key, msg)
	r.client.LTrim(ctx, key, -10, -1) // Keep only last 10 message
}

func (r *RedisStorage) GetHistory(userID int64) []string {
	ctx := context.Background()
	key := fmt.Sprintf("chat:%d", userID)
	message, _ := r.client.LRange(ctx, key, 0, -1).Result()
	return message
}
