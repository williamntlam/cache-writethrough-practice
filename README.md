# Go Practice: Writethrough Caching System

A simple task management application built with **Go (Gin)**, **Angular**, and **Redis** to practice implementing a writethrough caching system.

## 🎯 Learning Objectives

- **Go/Golang**: Web API development with Gin framework
- **Angular**: Frontend development and HTTP client usage
- **Caching Systems**: Understanding writethrough cache patterns
- **Database Design**: PostgreSQL with Redis caching layer
- **Docker**: Containerized development environment

## 🏗️ Architecture Overview

```
┌─────────────┐    HTTP    ┌─────────────┐    SQL    ┌─────────────┐
│   Angular   │ ────────── │   Go API    │ ───────── │ PostgreSQL  │
│  Frontend   │            │  (Gin)      │           │  Database   │
└─────────────┘            └─────────────┘           └─────────────┘
                                   │
                                   │ Redis Protocol
                                   ▼
                            ┌─────────────┐
                            │    Redis    │
                            │   Cache     │
                            └─────────────┘
```

### Writethrough Caching Pattern

1. **Write Operations**: Data is written to both the database AND cache simultaneously
2. **Read Operations**: Data is first checked in cache, then database if cache miss
3. **Consistency**: Ensures cache and database are always in sync

## 🚀 Quick Start

### Prerequisites

- Go 1.22+
- Docker & Docker Compose
- Node.js 18+ (for Angular)
- Git

### 1. Clone and Setup

```bash
git clone https://github.com/williamntlam/cache-writethrough-practice
cd cache-writethrough-practice
```

### 2. Start Infrastructure

```bash
# Start PostgreSQL and Redis containers
docker-compose -f docker/docker-compose.yml up -d
```

### 3. Setup Environment Variables

Create the required `.env` files:

**postgres/.env:**
```env
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=database
```

**redis/.env:**
```env
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

### 4. Run Go Backend

```bash
# Install dependencies
go mod tidy

# Run with hot reload (using Air)
./run-dev.sh

# Or run directly
go run backend/main.go
```

The API will be available at `http://localhost:8080`

### 5. Run Angular Frontend

```bash
cd frontend
npm install
ng serve
```

The frontend will be available at `http://localhost:4200`

## 📚 API Endpoints

### Tasks Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/tasks` | Get all tasks |
| `POST` | `/tasks` | Create a new task |
| `PATCH` | `/tasks/{id}` | Update a task |
| `DELETE` | `/tasks/{id}` | Delete a task |

### Example Usage

```bash
# Get all tasks
curl http://localhost:8080/tasks

# Create a task
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Go"}'

# Update a task
curl -X PATCH http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Go and Redis"}'

# Delete a task
curl -X DELETE http://localhost:8080/tasks/1
```

## 🗄️ Database Schema

### PostgreSQL Table

```sql
CREATE TABLE Tasks (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL
);
```

### Redis Cache Keys

- `task:{id}` - Individual task cache
- `tasks:all` - All tasks list cache

## 🔧 Project Structure

```
cache-writethrough-practice/
├── backend/
│   └── main.go              # Main Go application
├── frontend/                # Angular application
├── postgres/
│   ├── postgres.go          # Database connection
│   └── .env                 # Database config
├── redis/
│   ├── redis.go             # Redis connection
│   └── .env                 # Redis config
├── types/
│   └── tasks.go             # Go structs/types
├── docker/
│   ├── docker-compose.yml   # Infrastructure setup
│   ├── postgres/
│   │   ├── Dockerfile.postgres
│   │   └── init.sql
│   └── redis/
│       └── Dockerfile.redis
├── .air.toml               # Hot reload config
├── run-dev.sh              # Development script
└── README.md
```

## 🧪 Caching Implementation

### Writethrough Pattern

```go
// Write operation (POST/PATCH/DELETE)
func createTask(c *gin.Context) {
    // 1. Write to database
    err := db.QueryRow(`INSERT INTO Tasks(title) VALUES($1) RETURNING id`, title).Scan(&taskID)
    
    // 2. Write to cache (writethrough)
    cacheKey := fmt.Sprintf("task:%d", taskID)
    redisClient.Set(ctx, cacheKey, task, time.Hour)
    
    // 3. Invalidate list cache
    redisClient.Del(ctx, "tasks:all")
}

// Read operation (GET)
func getTasks(c *gin.Context) {
    // 1. Try cache first
    cached, err := redisClient.Get(ctx, "tasks:all").Result()
    if err == nil {
        // Cache hit - return cached data
        c.JSON(200, cached)
        return
    }
    
    // 2. Cache miss - query database
    rows, err := db.Query("SELECT id, title FROM Tasks")
    
    // 3. Update cache
    redisClient.Set(ctx, "tasks:all", tasks, time.Hour)
    
    c.JSON(200, tasks)
}
```

## 🎓 Learning Concepts

### Go/Golang
- **Gin Framework**: HTTP routing and middleware
- **Database/SQL**: PostgreSQL integration
- **Redis Client**: Caching operations
- **Environment Variables**: Configuration management
- **Error Handling**: Go error patterns

### Caching Systems
- **Writethrough Pattern**: Write to both cache and database
- **Cache Invalidation**: Managing stale data
- **Cache Keys**: Naming conventions
- **TTL (Time To Live)**: Cache expiration

### System Design
- **Microservices**: API separation
- **Containerization**: Docker setup
- **Environment Management**: Development vs production
- **API Design**: RESTful endpoints

## 🐛 Common Issues & Solutions

### Database Connection Issues
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check database logs
docker logs postgres
```

### Redis Connection Issues
```bash
# Check if Redis is running
docker ps | grep redis

# Test Redis connection
docker exec -it redis redis-cli ping
```

### Go Application Issues
```bash
# Check if application is running
ps aux | grep main

# Check application logs
# (logs should appear in the terminal where you ran the app)
```

## 🚀 Next Steps

1. **Implement Cache Invalidation**: Add proper cache invalidation strategies
2. **Add Authentication**: Implement JWT-based authentication
3. **Error Handling**: Improve error handling and logging
4. **Testing**: Add unit and integration tests
5. **Monitoring**: Add metrics and health checks
6. **Performance**: Implement connection pooling and optimization

## 📖 Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [Redis Documentation](https://redis.io/documentation)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Angular Documentation](https://angular.io/docs)

## 🤝 Contributing

This is a learning project. Feel free to:
- Add new features
- Improve documentation
- Fix bugs
- Suggest improvements

## 📄 License

This project is for educational purposes.
