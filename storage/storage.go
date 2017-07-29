package storage

import "github.com/go-redis/redis"

import "fmt"

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

// RedisErrorHandle handle redis error
func RedisErrorHandle(err error, exit bool, msg string) {
	fmt.Printf("Redis %s\nError: %s", msg, err)
	if exit {
		panic(err)
	}
}
