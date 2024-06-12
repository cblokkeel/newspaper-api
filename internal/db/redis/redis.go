package redisdb

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	RDB *redis.Client
}

func NewRedisDB() *RedisDB {
	RDB := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if err := RDB.Set(context.Background(), "foo", "bar", 0).Err(); err != nil {
		panic(err)
	}

	return &RedisDB{
		RDB,
	}
}
