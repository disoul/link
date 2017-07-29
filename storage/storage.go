package storage

import (
	"fmt"
	"link/model"

	"github.com/go-redis/redis"
)

// RedisConnection use connection struct to connect redis server
type RedisConnection struct {
	addr     string
	password string
	dbIndex  int
}

// DefaultConnection localhost:6379 default local redis server
var DefaultConnection = RedisConnection{"localhost:6379", "", 1}

// CreateClient : create and return a redis client
func CreateClient(conn RedisConnection) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     conn.addr,
		Password: conn.password,
		DB:       conn.dbIndex,
	})

	if _, err := client.Ping().Result(); err != nil {
		RedisErrorHandle(err, true, "send PING fail")
	}

	return client
}

// UpdateModel save/update model to redis
func UpdateModel(model model.Model, client *redis.Client) error {
	modelString, err := model.Stringify()
	if err != nil {
		return err
	}

	err = client.Set(
		fmt.Sprintf("link-model@%s", model.ID),
		modelString,
		0,
	).Err()
	if err != nil {
		return err
	}

	return nil
}

// RedisErrorHandle handle redis error
func RedisErrorHandle(err error, exit bool, msg string) {
	fmt.Printf("Redis %s\nError: %s", msg, err)
	if exit {
		panic(err)
	}
}
