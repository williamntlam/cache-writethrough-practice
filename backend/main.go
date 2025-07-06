package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type Todo struct {
	ID string `json:"id"`
	Title string `json:"title`
}

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