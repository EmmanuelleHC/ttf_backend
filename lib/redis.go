package lib

import (
	"context"
	"fmt"
	"time"

	"github.com/Aguztinus/petty-cash-backend/constants"
	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	cache  *cache.Cache
	client *redis.Client
	prefix string
}

// NewRedis creates a new redis client instance
func NewRedis(config Config, logger Logger) Redis {
	addr := config.Redis.Addr()

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       constants.RedisMainDB,
		Password: config.Redis.Password,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		logger.Zap.Fatalf("Error to open redis[%s] connection: %v", addr, err)
	}

	logger.Zap.Info("Redis connection established")
	return Redis{
		client: client,
		prefix: config.Redis.KeyPrefix,
		cache: cache.New(&cache.Options{
			Redis:      client,
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		}),
	}
}

func (a Redis) wrapperKey(key string) string {
	return fmt.Sprintf("%s:%s", a.prefix, key)
}

func (a Redis) Set(key string, value interface{}, expiration time.Duration) error {
	return a.cache.Set(&cache.Item{
		Ctx:            context.TODO(),
		Key:            a.wrapperKey(key),
		Value:          value,
		TTL:            expiration,
		SkipLocalCache: true,
	})
}

func (a Redis) Get(key string, value interface{}) error {
	err := a.cache.Get(context.TODO(), a.wrapperKey(key), value)
	if err == cache.ErrCacheMiss {
		err = errors.RedisKeyNoExist
	}

	return err
}

func (a Redis) Do(key string, value interface{}) (string, error) {
	val, err := a.client.Do(context.TODO(), a.wrapperKey(key), value).Text()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("key does not exists")
			return "", err
		}
	}

	return val, nil
}

func (a Redis) Delete(keys ...string) (bool, error) {
	wrapperKeys := make([]string, len(keys))
	for index, key := range keys {
		wrapperKeys[index] = a.wrapperKey(key)
	}

	cmd := a.client.Del(context.TODO(), wrapperKeys...)

	if err := cmd.Err(); err != nil {
		return false, err
	}

	return cmd.Val() > 0, nil
}

func (a Redis) Increment(key string) (int64, error) {
	cmd := a.client.Incr(context.TODO(), key)
	if err := cmd.Err(); err != nil {
		return 0, err
	}

	return cmd.Val(), nil
}

func (a Redis) Check(keys ...string) (bool, error) {
	wrapperKeys := make([]string, len(keys))
	for index, key := range keys {
		wrapperKeys[index] = a.wrapperKey(key)
	}

	cmd := a.client.Exists(context.TODO(), wrapperKeys...)
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

// SAdd call redis SADD function
func (a Redis) SAdd(ctx context.Context, key, member string) error {
	err := a.client.SAdd(ctx, key, member).Err()
	if err != nil {
		return err
	}
	return nil
}

// SMembers return all members in a set
func (a Redis) SMembers(ctx context.Context, key string) ([]string, error) {
	result, err := a.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// SRem call redis SREM function
func (a Redis) SRem(ctx context.Context, key string, members ...string) error {
	err := a.client.SRem(ctx, key, members).Err()
	if err != nil {
		return err
	}
	return nil
}

func (a Redis) Close() error {
	return a.client.Close()
}

func (a Redis) GetClient() *redis.Client {
	return a.client
}
