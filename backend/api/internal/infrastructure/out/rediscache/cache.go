package rediscache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/lregnier/design-youtube/api/internal/application"
)

var _ application.Cache = (*Cache)(nil)

const ttl = 60 * time.Second

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, errors.New("cache miss")
	}
	return val, err
}

func (c *Cache) Set(ctx context.Context, key string, value []byte) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}
