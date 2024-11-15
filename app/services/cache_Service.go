package services

import (
	"context"
	"fmt"
	"go-microservice-product-porto/app/models"
	"log"
	"time"

	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type CacheService struct {
	client     *redis.Client
	defaultTTL time.Duration
	listTTL    time.Duration
}

func NewCacheService(client *redis.Client) *CacheService {
	return &CacheService{
		client:     client,
		defaultTTL: 24 * time.Hour,  // Default TTL for single products
		listTTL:    5 * time.Minute, // Default TTL for product lists
	}
}

func (s *CacheService) SetDefaultTTL(duration time.Duration) {
	s.defaultTTL = duration
}

func (s *CacheService) SetListTTL(duration time.Duration) {
	s.listTTL = duration
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

	return s.client.Set(context.Background(), key, jsonData, s.defaultTTL).Err()
}

func (s *CacheService) Delete(key string) error {
	return s.client.Del(context.Background(), key).Err()
}

func (s *CacheService) GetList(key string) (*models.PaginatedResponse, error) {
	val, err := s.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	var response models.PaginatedResponse
	err = json.Unmarshal([]byte(val), &response)
	return &response, err
}

func (s *CacheService) SetList(key string, value *models.PaginatedResponse) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.client.Set(context.Background(), key, json, s.listTTL).Err()
}

func (s *CacheService) DeletePattern(pattern string) error {
	iter := s.client.Scan(context.Background(), 0, pattern, 0).Iterator()
	for iter.Next(context.Background()) {
		err := s.client.Del(context.Background(), iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	return iter.Err()
}

func (s *CacheService) GenerateProductKey(id string) string {
	return fmt.Sprintf("product:%s", id)
}

func (s *CacheService) GenerateListKey(page int, pageSize int, sortBy string, sortDir string) string {
	return fmt.Sprintf("products:list:p%d:s%d:sort_%s_%s", page, pageSize, sortBy, sortDir)
}

func (s *CacheService) InvalidateRelatedCaches(productID string) error {
	err := s.Delete(s.GenerateProductKey(productID))
	if err != nil {
		return err
	}

	return s.DeletePattern("products:list:*")
}
