package redis

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	// Client is the Redis client
	Client *redis.Client
	// Ctx is the context for Redis operations
	Ctx = context.Background()
)

// InitRedis initializes the Redis connection
func InitRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0

	// Parse DB number if provided
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		var err error
		redisDB, err = strconv.Atoi(dbStr)
		if err != nil {
			log.Printf("Invalid REDIS_DB value, using default 0: %v", err)
			redisDB = 0
		}
	}

	Client = redis.NewClient(&redis.Options{
		Addr:         redisHost + ":" + redisPort,
		Password:     redisPassword,
		DB:           redisDB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	})

	// Test the connection
	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
}

// Set stores a value in Redis with expiration
func Set(key string, value interface{}, expiration time.Duration) error {
	return Client.Set(Ctx, key, value, expiration).Err()
}

// Get retrieves a value from Redis
func Get(key string) (string, error) {
	return Client.Get(Ctx, key).Result()
}

// Delete removes a key from Redis
func Delete(key string) error {
	return Client.Del(Ctx, key).Err()
}

// CloseRedis closes the Redis connection
func CloseRedis() {
	if Client != nil {
		Client.Close()
	}
}