package main

import (
	"context"
	"github.com/gin-gonic/gin"
	// "github.com/redis/go-redis/v9"
	// "fmt"
	// "net/http"
	// "cache-writethrough-practice/postgres"
	// "cache-writethrough-practice/redis"
)

var ctx = context.Background()

func main() {

	router := gin.Default()

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