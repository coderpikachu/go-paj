package my_cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"time"
)

func MyCacheTest() string {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "101.43.168.151:6379",
		Password: "1436001", // no password set
		DB:       0,         // use default DB
	})
	pong, err := rdb.Ping(context.Background()).Result()
	fmt.Println(pong, err)
	c := NewRedisCache(rdb)
	err = c.Set(context.Background(), "key1", "abc", time.Minute)
	if err != nil {
		return "30*"
	}
	val, err := c.Get(context.Background(), "key1")
	return val.(string)
}
