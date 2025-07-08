package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
	redisdb "github.com/redis/go-redis/v9"
	"cache-writethrough-practice/postgres"
	"cache-writethrough-practice/redis"
	"cache-writethrough-practice/types"
)

var ctx = context.Background()
var db *sql.DB
var redisClient *redisdb.Client

func main() {

	router := gin.Default()

	var err error
	db, err = postgres.ConnectToPostgres()

	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	} 

	defer db.Close()
	fmt.Println("Connected to Postgres");

	redisClient, err = redis.ConnectToRedis()

	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	defer redisClient.Close()
	fmt.Println("Connected to Redis")

	// POST - create a task - uses cache
	// GET - get all tasks. - uses cache
	// GET - get a specific task - uses cache
	// PATCH - update a task
	// DELETE - delete a task

	router.PATCH("/tasks/:id", func(context *gin.Context) {
		taskIDStr := context.Param("id")

		var request types.TaskRequest
		
		var taskID int
		if _, err := fmt.Sscanf(taskIDStr, "%d", &taskID); err != nil {
			context.JSON(400, gin.H{"error": "Invalid task ID"})
			return
		}
	
		if err := context.ShouldBindJSON(&request); err != nil {
			context.JSON(400, gin.H{"error": err.Error()})
			return
		}
	
		result, err := db.Exec(`UPDATE Tasks SET title = $1 WHERE id = $2`, request.Title, taskID)
		if err != nil {
			context.JSON(500, gin.H{"error": "Failed to update task"})
			return
		}
	
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			context.JSON(500, gin.H{"error": "Failed to update task"})
			return
		}
	
		if rowsAffected == 0 {
			context.JSON(404, gin.H{"error": "Task not found"})
			return
		}
	
		key := fmt.Sprintf("task:%d", taskID)
		err = redisClient.HSet(ctx, key, "title", request.Title).Err()

		if err != nil {
			context.JSON(404, gin.H{"error": "Task could not be updated in the Redis cache."})
			return
		}

		context.JSON(200, gin.H{"message": "Task updated successfully"})
	})

	router.DELETE("/tasks/:id", func(context *gin.Context) {

		taskIDStr := context.Param("id");
		
		var taskID int
		if _, err := fmt.Sscanf(taskIDStr, "%d", &taskID); err != nil {
			context.JSON(400, gin.H{"error": "Invalid task ID"})
			return
		}

		result, err := db.Exec(`DELETE from Tasks where id = $1`, taskID)

		if err != nil {
			fmt.Println("Delete Failed: %v", err)
			context.JSON(500, gin.H{"error": "Failed to delete task"})
			return
		}

		rowsAffected, err := result.RowsAffected()

		if err != nil {
			log.Printf("Error getting rows affected: %v", err)
			context.JSON(500, gin.H{"error": "Failed to delete task"})
			return
		}

		if rowsAffected == 0 {
			context.JSON(404, gin.H{"error": "Task not found"})
			return
		}

		// Delete task from the Redis cache.

		key := fmt.Sprintf("task:%d", taskID)
		err = redisClient.Del(ctx, key).Err()

		if err != nil {
			context.JSON(400, gin.H{"error": "Task not found in the Redis Cache"})
			return
		}

		context.JSON(200, gin.H{"message": "Task deleted successfully"})

	})

	router.GET("/tasks/:id", func(context *gin.Context) {
		taskIDStr := context.Param("id")
		var taskID int
		if _, err := fmt.Sscanf(taskIDStr, "%d", &taskID); err != nil {
			context.JSON(400, gin.H{"error": "Invalid task ID"})
			return
		}

		key := fmt.Sprintf("task:%d", taskID)

		result, err := redisClient.HGetAll(ctx, key).Result()

		if err != nil {
			context.JSON(400, gin.H{"error": "Redis Error"})
			return
		} 
	
		if len(result) == 0 {

			context.JSON(404, gin.H{"error": "Task not found"})
			return

		}

		id := result["id"]
		title := result["title"]

		context.JSON(200, gin.H{
			"id":    id,
			"title": title,
		})
	})

	router.GET("/tasks", func(context *gin.Context) {

		keys, err := redisClient.Keys(ctx, "task:*").Result()

		if err != nil {
			log.Printf("Query failed: %v", err)
			context.JSON(500, gin.H{"error": "Failed to fetch tasks"})
			return
		}

		tasks := make([]gin.H, 0)
		for _, key := range keys {
			task, err := redisClient.HGetAll(ctx, key).Result()

			if err != nil {
				log.Printf("Error processing tasks: %v", err)
				context.JSON(500, gin.H{"error": "Failed to process tasks"})
				return
			}

			if len(task) > 0 {
				tasks = append(tasks, gin.H(task))
			}

		}

		context.JSON(200, gin.H{
			"tasks": tasks,
		})
	})

	router.POST("/tasks", func(context *gin.Context) {

		var request types.TaskRequest

		if err := context.ShouldBindJSON(&request); err != nil {
			context.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var taskID int

		// write to the postgres database next.
		err := db.QueryRow(`INSERT INTO Tasks(title) VALUES($1) RETURNING id`, request.Title).Scan(&taskID)

		if err != nil {
			context.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// write to the redis cache next.
		// request.Title
		// taskID

		err = redisClient.HSet(ctx, fmt.Sprintf("task:%d", taskID), map[string]interface{}{
			"id": taskID,
			"title": request.Title,
		}).Err()

		if err != nil {
			context.JSON(500, gin.H{"error": err.Error()})
			return
		}

		context.JSON(201, gin.H{
			"message": "Task created successfully",
			"task_id": taskID,
		})

	})

	router.Run()

}