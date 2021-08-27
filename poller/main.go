package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis"
)

var (
	redisConn *redis.Client
)

func init() {
	redisConn = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		// DB:       0,  // use default DB
	})
}

func Poll(interval time.Duration, done <-chan os.Signal) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			fmt.Println("scheduled tasks at : ", time.Now())
			now := time.Now().Unix()
			res := redisConn.Do("ZRANGEBYSCORE", "key:task", "-inf", float64(now)).Val()

			s := res.([]interface{})
			if len(s) == 0 {
				fmt.Println("task is empty")
			}

			for _, task := range s {
				fmt.Println("task : ", task)

				// do something task...

				// Delete task from redis
				_, err := redisConn.ZRem("key:task", task).Result()
				if err != nil {
					fmt.Errorf("delete value in redis error : %v", err)
				}
				fmt.Println("task is success")
			}
		case <-done:
			fmt.Println("shutting down poller")
			return
		}
	}
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	fmt.Println("scheduled tasks running every 15 second")
	Poll(15*time.Second, c)
}
