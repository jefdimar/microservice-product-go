package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"go-microservice-product-porto/pkg/errors"
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
		return nil, errors.StandardError(errors.ECACHE, fmt.Errorf("failed to connect to Redis: %v", err))
	}

	return &RedisCache{client: client}, nil
}

func (c *RedisCache) Set(key string, value interface{}) error {
	json, err := json.Marshal(value)
	if err != nil {
		return errors.StandardError(errors.ECACHE, fmt.Errorf("failed to marshal value: %v", err))
	}

	err = c.client.Set(context.Background(), key, json, 24*time.Hour).Err()
	if err != nil {
		return errors.StandardError(errors.ECACHE, fmt.Errorf("failed to set value in Redis: %v", err))
	}

	return nil
}

func (c *RedisCache) Get(key string) (interface{}, error) {
	val, err := c.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, errors.StandardError(errors.ENOTFOUND, fmt.Errorf("key not found in Redis: %v", err))
	}
	if err != nil {
		return nil, errors.StandardError(errors.ECACHE, fmt.Errorf("failed to get value from Redis: %v", err))
	}

	var result interface{}
	if err = json.Unmarshal([]byte(val), &result); err != nil {
		return nil, errors.StandardError(errors.ECACHE, fmt.Errorf("failed to unmarshal value from Redis: %v", err))
	}

	return result, nil
}

func (c *RedisCache) Delete(key string) error {
	err := c.client.Del(context.Background(), key).Err()
	if err != nil {
		return errors.StandardError(errors.ECACHE, fmt.Errorf("failed to delete key from Redis: %v", err))
	}

	return nil
}
