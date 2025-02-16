package cache

import "go-microservice-product-porto/internal/infrastructure/persistence/redis"

type CacheService interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
}

func NewCacheService(redisConfig redis.RedisConfig) (CacheService, error) {
	return redis.NewRedisCache(redisConfig)
}
