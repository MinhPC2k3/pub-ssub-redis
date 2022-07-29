package main

import (
	"context"
	"fmt"
	"time"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	channelName = "mychannel"
	ctx = context.Background()
)

func InitRedisClient () *redis.Client{
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}
func main() {
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "",
	// 	DB:       0,
	// })
	rdb:=InitRedisClient()
	var wg sync.WaitGroup
	wg.Add(2)
	go sendData(rdb)
	go printData(rdb)
	wg.Wait()
}
func sendData(rdb *redis.Client) {
	num := 0
	for {
		err := rdb.Publish(ctx, channelName, num).Err()
		if err != nil {
			panic(err)
		}
		// fmt.Println("process1: ", num)
		num++
		time.Sleep(time.Second * 1)
	}
}
func printData(rdb *redis.Client) {
	sub := rdb.Subscribe(ctx, channelName)
	c := 0
	for {
		message, err := sub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println(c, message)
		c++
	}
}
