package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCacheService struct {
	ctx context.Context
	rdb *redis.Client
}

func NewRedisCacheService(rdb *redis.Client) RedisCacheService {
	return &redisCacheService{
		ctx: context.Background(),
		rdb: rdb,
	}
}

// Get retrieves a cached value from Redis by the specified key.
// The data is unmarshaled from JSON into the provided destination struct.
func (cs *redisCacheService) Get(key string, dest any) error {
	data, err := cs.rdb.Get(cs.ctx, key).Result()

	if err == redis.Nil {
		return err
	}

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Set stores a given value in Redis under the specified key with a TTL (Time-To-Live).
// The value is marshaled into JSON before being saved.
func (cs *redisCacheService) Set(key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return cs.rdb.Set(cs.ctx, key, data, ttl).Err()
}

// Clear deletes all keys from Redis that match the given pattern.
// It uses the SCAN command to iterate through keys safely without blocking Redis.
func (cs *redisCacheService) Clear(pattern string) error {
	cursor := uint64(0)
	for {
		keys, nextCursor, err := cs.rdb.Scan(cs.ctx, cursor, pattern, 2).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			cs.rdb.Del(cs.ctx, keys...)
		}

		cursor = nextCursor

		if cursor == 0 {
			break
		}
	}

	return nil
}

// Exists checks whether a specific key exists in Redis.
func (cs *redisCacheService) Exists(key string) (bool, error) {
	count, err := cs.rdb.Exists(cs.ctx, key).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
