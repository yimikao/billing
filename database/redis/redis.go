package redisDB

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/yimikao/billing/core"
)

type CacheKey string

const (
	Userdata = CacheKey("userdata")
)

type Cache interface {
	StoreData(context.Context, CacheKey, interface{}) error
	GetData(context.Context, CacheKey) (string, error)
}

type Client struct {
	inner *redis.Client
}

func New(dsn string) (*Client, error) {

	opts, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, err
	}

	opts.Password = core.Global().Database.Redis.Password

	if core.Global().Database.Redis.UseTLS {
		opts.TLSConfig = &tls.Config{}
	}

	if strings.Contains(dsn, "localhost") || strings.Contains(dsn, "127.0.0.1") ||
		strings.Contains(dsn, "host.docker.internal") {
		opts.TLSConfig = nil
	}

	client := redis.NewClient(opts)

	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*15)
	defer cancelFn()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Client{
		inner: client,
	}, nil
}

func (c *Client) StoreData(ctx context.Context, key CacheKey, data interface{}) error {
	return c.inner.Set(ctx, string(key), data, 0).Err()
}

func (c *Client) GetData(ctx context.Context, key CacheKey) (string, error) {
	return c.inner.Get(ctx, string(key)).Result()
}
