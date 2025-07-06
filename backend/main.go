package main

import (
	"context"
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
	"cache-writethrough-practice/postgres"
	"cache-writethrough-practice/redis"
)

var ctx = context.Background()

func main() {

	router := gin.Default()

	db, err := postgres.ConnectToPostgres()

	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
		os.Exit(1);
	} 

	defer db.Close()
	fmt.Println("Connected to Postgres");

	redisClient, err := redis.ConnectToRedis()

	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
		os.Exit(1);
	}

	defer redisClient.Close()
	fmt.Println("Connected to Redis")

	// POST - create a task
	// GET - get a task
	// GET - get all tasks.
	// PUT/PATCH - update a task
	// DELETE - delete a task

	router.GET("/ping", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run()

}