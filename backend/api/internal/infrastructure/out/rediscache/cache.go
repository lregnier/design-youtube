package rediscache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/lregnier/design-youtube/api/internal/application"
)

var _ application.Cache = (*cache)(nil)

const ttl = 60 * time.Second

type cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) application.Cache {
	return &cache{client: client}
}

func (c *cache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, errors.New("cache miss")
	}
	return val, err
}

func (c *cache) Set(ctx context.Context, key string, value []byte) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}
