package mocks

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	redisDB "github.com/yimikao/billing/database/redis"
)

type MockCache struct {
	Inner *redis.Client
	Mock  redismock.ClientMock
}

func NewMockCache() *MockCache {
	redis, mock := redismock.NewClientMock()

	return &MockCache{
		Inner: redis,
		Mock:  mock,
	}
}

func (c *MockCache) StoreData(ctx context.Context, key redisDB.CacheKey, data interface{}) error {
	// return c.Inner.Set(ctx, string(key), data, 0).Err()
	return nil
}

func (c *MockCache) GetData(ctx context.Context, key redisDB.CacheKey) (string, error) {
	// return c.Inner.Get(ctx, string(key)).Result()
	return "", nil
}
