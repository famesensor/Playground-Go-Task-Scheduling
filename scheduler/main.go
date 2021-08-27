package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var (
	redisConn *redis.Client
)

type Task struct {
	Id     string
	Detail string
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
	saveTaskToRedis("key:task", Task{Id: "task_uuid_1", Detail: "You must be learn everything", Status: "waiting"}, time.Now().Add(time.Second*10))
	saveTaskToRedis("key:task", Task{Id: "task_uuid_2", Detail: "You must be go to shop", Status: "waiting"}, time.Now().Add(time.Minute*15))
}

func saveTaskToRedis(keyGroup string, task Task, date time.Time) {
	data := buildModelToJson(task)

	_, err := redisConn.ZAdd(keyGroup, redis.Z{Score: float64(date.Unix()), Member: data}).Result()
	if err != nil {
		fmt.Errorf("add to redis is error : %v", err)
	}

	fmt.Printf("Add task to redis success id : %v, at : %v", task.Id, date)
}

func buildModelToJson(value interface{}) []byte {
	js, err := json.Marshal(value)
	if err != nil {
		fmt.Errorf("build model to json error : %v", err)
	}

	return js
}
