package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	redisConn *redis.Client
	ctx       = context.Background()
)

type DateModel struct {
	Id     string
	Status string
}

func init() {
	redisConn = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		// DB:       0,  // use default DB
	})
}

func main() {
	md := DateModel{
		Status: "expire",
	}
	saveTaskToRedis("promotion", md, time.Now().Add(time.Second*10))
}

func saveTaskToRedis(key string, dateModel DateModel, date time.Time) {
	fmt.Printf("key : %v, status : %v, date : %v", key, dateModel.Status, date.Unix())
	err := redisConn.Do(ctx, "ZADD", key, dateModel, float64(date.Unix())).Err()
	if err != nil {
		fmt.Errorf("save to redis error : %v", err)
	}
}
