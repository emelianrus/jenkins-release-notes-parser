package main

import (
	"errors"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func NewRedisClient() *Redis {
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		// Password: "",
		DB: 0,
	})
	return &Redis{client: client}
}

func (r *Redis) Get(key string) *redis.StringCmd {
	return r.client.Get(key)
}

func (r *Redis) Set(key string, value interface{}) error {
	return r.client.Set(key, value, 0).Err()
}

func (r *Redis) GetJenkinsServers() ([]byte, error) {
	return r.client.Get("servers:jenkins-one:plugins").Bytes()
}

func (r *Redis) GetPlugin(key string) ([]byte, error) {
	jsonData, err := r.Get(key).Bytes()
	if err != nil {
		return []byte{}, errors.New("error in getPlugins " + key)
	}
	return jsonData, err
}

// func InitDB() *redis.Client {
// 	log.Println("Redis client init")
// 	client := redis.NewClient(&redis.Options{
// 		Addr: "127.0.0.1:6379",
// 		// Password: "",
// 		DB: 0,
// 	})
// 	return client
// }

// func GetByName(){

// }

// func Set(){

// 	err = redisclient.Set(fmt.Sprintf("github:%s:%s:%s", ownerName, repoName, "lastUpdated"),
// 		time.Now().Unix(), 0).Err()
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// }

// func Get(){

// }
// jsonData, err := redisclient.Get(key).Bytes()
// 		if err != nil {
// 			log.Println(err)
// 			http.Error(w, "Failed to retrieve releases from cache", http.StatusInternalServerError)
// 			return
// 		}

// err = redisclient.Set(fmt.Sprintf("github:%s:%s:%s", ownerName, repoName, "lastUpdated"),
// 		time.Now().Unix(), 0).Err()

// 		\
