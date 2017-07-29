package main

import (
	"main/storage"
	"time"
)

func initRedis() *redis.Client {
	redisClient := storage.CreateClient(storage.DefaultConnection)

	_, err := redisClient.Set("link@uptime", time.Now().Unix(), 0).Result()
	if err != nil {
		storage.RedisErrorHandle(err, true, "send SET command fail, canont set link@uptime")
	}

	return redisClient
}

func main() {
	initRedis()
}
