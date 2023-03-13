package main

import (
	"log"

	"github.com/go-redis/redis"
)

func InitDB() *redis.Client {
	log.Println("Redis client init")
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		// Password: "",
		DB: 0,
	})
	return client
}
