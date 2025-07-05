import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
	"github.com/joho/godotenv"
)

var RedisDatabase *redis.Client

func connectToRedis() {

	godotenv.Load()
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	RedisDatabase = redis.NewClient(&redis.Options{

		Addr:     fmt.Sprintf("%s:%s", host, port),
        Password: "", // or os.Getenv("REDIS_PASSWORD")
        DB:       0,  // default DB

	})

	_, error = RedisDatabase.Ping(context.Background()).Result()

	if error != nil {

		panic("Failed to connect to Redis: " + error.Error())

	}

	fmt.Println("Connected to Redis")

}