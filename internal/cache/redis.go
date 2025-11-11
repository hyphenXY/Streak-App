// cache/redis.go
package cache

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	Rdb *redis.Client
)

func InitRedis() {
	// Load Redis config from environment with sensible defaults
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		host := os.Getenv("REDIS_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("REDIS_PORT")
		if port == "" {
			port = "6379"
		}
		addr = host + ":" + port
	}

	password := os.Getenv("REDIS_PASSWORD") // empty = no password

	db := 0
	if v := os.Getenv("REDIS_DB"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			db = n
		} else {
			log.Printf("invalid REDIS_DB=%q, defaulting to DB 0", v)
		}
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	fmt.Println("✅ Redis connected successfully!")
}

func GetCache(key string) (string, error) {
	val, err := Rdb.Get(Ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func SetCache(key string, value string, expirationSeconds int) error {
	err := Rdb.Set(Ctx, key, value, 0).Err()
	return err
}

func DeleteCache(key string) error {
	err := Rdb.Del(Ctx, key).Err()
	return err
}

func GetFromCacheOrDB(key string, dbFetchFunc func() (string, error), expirationSeconds int) (string, error) {
	// Try to get from cache
	val, err := GetCache(key)
	if err == nil {
		return val, nil // Cache hit
	}

	// Cache miss, fetch from DB
	val, err = dbFetchFunc()
	if err != nil {
		return "", err
	}

	// Set the fetched value in cache
	err = SetCache(key, val, expirationSeconds)
	if err != nil {
		return "", err
	}

	return val, nil
}
