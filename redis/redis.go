package redisdb

import (
    "context"
    "fmt"
    "os"

    "github.com/joho/godotenv"
    "github.com/redis/go-redis/v9"
)

func ConnectToRedis() (*redis.Client, error) {
    // Load environment variables
    err := godotenv.Load("./redis/.env")

	if err != nil {
		fmt.Println("Warning: No .env file found in ./postgres/.env")
	}

    host := os.Getenv("REDIS_HOST")
    port := os.Getenv("REDIS_PORT")

    rdb := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%s", host, port),
        Password: os.Getenv("REDIS_PASSWORD"), // optional
        DB:       0, // default DB
    })

    // Test the connection
    pong, err := rdb.Ping(context.Background()).Result()
    if err != nil {
        return nil, fmt.Errorf("failed to connect to Redis: %w", err)
    }

    fmt.Println("Redis connected:", pong)
    return rdb, nil
}