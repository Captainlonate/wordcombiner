package conceptsdb

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

// Package variable that will hold the Redis connection
var rdb *redis.Client

func InitializeRedisConnection() {
	// Establish connection to Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Ping the Redis server to check if it's running
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}

	fmt.Println("Connected to Redis:", pong)
}

// Retrieve a value from redis by it's key.
// If the value is not found, an empty string will be returned.
// An error will be returned if there is an actual problem.
func FetchFromRedis(key string) (string, error) {
	retrievedValue, err := rdb.Get(context.Background(), key).Result()

	// If there is an actual error, not just "key not found", but a real error
	if err != nil && err != redis.Nil {
		return "", err
	}

	return retrievedValue, nil
}

// Adds an entry into Redis.
func InsertToRedis(key string, value string) error {
	err := rdb.Set(context.Background(), key, value, 0).Err()

	if err != nil {
		return err
	}

	return nil
}
