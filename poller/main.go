package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	redisConn *redis.Client
	ctx       = context.Background()
)

func init() {
	redisConn = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		// DB:       0,  // use default DB
	})
}

// Poll checks Redis to determine whether scheduled tasks need to be run or not.
func Poll(interval time.Duration, done <-chan os.Signal) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			fmt.Println("Checking for scheduled tasks... at : ", time.Now())
			response, err := redisConn.Do(ctx, "ZREMRANGEBYSCORE", "promotion", "-inf").Result()
			if err != nil {
				fmt.Errorf("error redis : %v", err)
			}
			fmt.Println(response)
			// if len(response) == 0 {
			// 	fmt.Println("task is empty")
			// }

			// for _, res := range response {
			// 	fmt.Printf("task is redis : %v\n", res)
			// }
		case <-done:
			fmt.Println("Shutting down poller")
			return
		}
	}
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	fmt.Println("Polling Redis every 10 seconds for scheduled tasks...")
	Poll(10*time.Second, c)
}
