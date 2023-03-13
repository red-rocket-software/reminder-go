package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/red-rocket-software/reminder-go/config"
	"github.com/redis/go-redis/v9"
)

type CacheRedis interface {
	Set(ctx context.Context, key string, val string) error
	Get(ctx context.Context, key string) (string, error)
	IfExistsInCache(ctx context.Context, key string) (bool, error)
}

type CacheConn struct {
	client     *redis.Client
	expiration time.Duration
}

func NewCacheConn(config config.Config) (*CacheConn, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addrs,
		Password: config.Redis.Password,
		DB:       0,
	})
	ctx := context.Background()
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		// Sleep for 3 seconds and wait for Redis to initialize
		time.Sleep(3 * time.Second)
		err := redisClient.Ping(ctx).Err()
		if err != nil {
			return nil, err
		}
	}
	fmt.Println(pong)

	return &CacheConn{
		client:     redisClient,
		expiration: time.Duration(config.Redis.ExpirationHour) * time.Second,
	}, nil
}

// Set sets a key-value pair
func (cache *CacheConn) Set(ctx context.Context, key string, val string) error {
	if err := cache.client.Set(ctx, key, val, cache.expiration).Err(); err != nil {
		return err
	}
	return nil
}

// Get returns true if the key already exists and set dst to the corresponding value
func (cache *CacheConn) Get(ctx context.Context, key string) (string, error) {
	val, err := cache.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (cache *CacheConn) IfExistsInCache(ctx context.Context, key string) (bool, error) {
	exist, err := cache.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if exist != 1 {
		return false, nil
	}
	return true, nil
}
