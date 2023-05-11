package common

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

type DataService struct {
  Client *redis.Client 
}

func NewDataService() (*DataService, error) {
	log.SetOutput(os.Stdout)

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisHost == "" {
		redisHost = "localhost"
	}
	if redisPort == "" {
		redisPort = "6379"
	}

	client := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%s", redisHost, redisPort), DB: 0})

  db := &DataService{Client: client}

	res, err := db.Client.Ping(context.Background()).Result()
  log.Println(res)
  return db, err
}

