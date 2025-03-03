package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	rdb redis.UniversalClient
}

type Cacher interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
}

func NewRedisCache(address, password string) (Cacher, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
	})
	ctx := context.Background()
	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	return &RedisCache{rdb: client}, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.rdb.Set(ctx, key, value, expiration).Err()
}

func (c *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	return c.rdb.Get(ctx, key).Bytes()
}
