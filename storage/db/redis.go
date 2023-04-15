package db

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Redis struct {
	client *redis.Client
}

var ctx = context.Background()

func NewRedisClient() *Redis {
	logrus.Infoln("Creating redis connection")
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
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
