package db

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Redis struct {
	client *redis.Client
}

var ctx = context.Background()

func getEnvOrDefault(envName string, defaultValue string) string {
	value := os.Getenv(envName)
	if value == "" {
		return defaultValue
	} else {
		return value
	}
}

func NewRedisClient() *Redis {
	logrus.Infoln("Creating redis connection")

	redisHost := getEnvOrDefault("REDIS_HOST", "127.0.0.1")
	redisPort := getEnvOrDefault("REDIS_PORT", "6379")

	client := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
		// Password: "",
		DB: 0,
	})

	return &Redis{client: client}
}

func (r *Redis) Status() error {
	return r.client.Ping(ctx).Err()
}

func (r *Redis) Get(key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

func (r *Redis) Set(key string, value interface{}) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *Redis) Keys(key string) ([]string, error) {
	return r.client.Keys(ctx, key).Result()
}

// TODO: not used
func (r *Redis) Del(key string) *redis.IntCmd {
	return r.client.Del(ctx, key)
}
