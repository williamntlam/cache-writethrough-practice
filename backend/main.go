package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
	"cache-writethrough-practice/postgres"
	"cache-writethrough-practice/redis"
)

var ctx = context.Background()

func main() {

	router := gin.Default()

	// POST - create a task
	// GET - get a task
	// GET - get all tasks.
	// PUT/PATCH - update a task
	// DELETE - delete a task

	

	router.Run()

}