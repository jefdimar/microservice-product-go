package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(cfg RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

func (c *RedisCache) Set(key string, value interface{}) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(context.Background(), key, json, 24*time.Hour).Err()
}

func (c *RedisCache) Get(key string) (interface{}, error) {
	val, err := c.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = json.Unmarshal([]byte(val), &result)
	return result, err
}

func (c *RedisCache) Delete(key string) error {
	return c.client.Del(context.Background(), key).Err()
}
