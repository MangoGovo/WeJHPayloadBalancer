package Dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var RedisDB *redis.Client

func InitRedis() {

	RedisDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis地址
		Password: "262626mjj",      // 密码（如果设置了的话）
		DB:       0,                // 使用默认DB
	})
}

func RedisSetKeyVal(ctx context.Context, key string, val string, expire time.Duration) error {
	return RedisDB.SetEX(ctx, key, val, expire).Err()
}

func RedisGetKeyVal(ctx context.Context, key string) (string, error) {
	return RedisDB.Get(ctx, key).Result()
}

func RedisDelKeyVal(ctx context.Context, key string) error {
	return RedisDB.Del(ctx, key).Err()
}
