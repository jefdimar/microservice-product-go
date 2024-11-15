package services

import (
	"context"
	"go-microservice-product-porto/app/models"
	"log"
	"time"

	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type CacheService struct {
	client *redis.Client
}

func NewCacheService(client *redis.Client) *CacheService {
	return &CacheService{
		client: client,
	}
}

func (s *CacheService) Get(key string) (*models.Product, error) {
	val, err := s.client.Get(context.Background(), key).Result()
	if err != nil {
		log.Printf("Cache error: %v", err)
		return nil, err
	}

	var product models.Product
	err = json.Unmarshal([]byte(val), &product)
	return &product, err
}

func (s *CacheService) Set(key string, product *models.Product) error {
	jsonData, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return s.client.Set(context.Background(), key, jsonData, 24*time.Hour).Err()
}

func (s *CacheService) Delete(key string) error {
	return s.client.Del(context.Background(), key).Err()
}

func (c *CacheService) GetList(key string) (*models.PaginatedResponse, error) {
	val, err := c.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	var response models.PaginatedResponse
	err = json.Unmarshal([]byte(val), &response)
	return &response, err
}

func (c *CacheService) SetList(key string, value *models.PaginatedResponse) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(context.Background(), key, json, 5*time.Minute).Err()
}

func (c *CacheService) DeletePattern(pattern string) error {
	iter := c.client.Scan(context.Background(), 0, pattern, 0).Iterator()
	for iter.Next(context.Background()) {
		err := c.client.Del(context.Background(), iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	return iter.Err()
}
