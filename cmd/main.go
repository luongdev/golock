package main

import (
	"context"
	"github.com/luongdev/golock/redis"
	libredis "github.com/redis/go-redis/v9"
	"log"
	"time"
)

func main() {
	store, err := redis.NewRedisLockStore("DefaultLockKey", libredis.NewClient(&libredis.Options{
		Addr:     "34.124.240.249:6979",
		Password: "",
		DB:       0,
	}))
	if err != nil {
		panic(err)
	}

	locker := redis.NewRedisLocker(store)

	for i := 0; i < 3; i++ {
		go func() {
			_, err = locker.Lock(
				redis.WithName("test"),
				redis.WithLockAtLeast(time.Second*5),
				redis.WithLockAtMost(time.Second*120),
			)
			if err != nil {
				log.Printf("error: %v", err)
			}
		}()
	}

	lock, err := store.Get(context.Background(), "test")
	if err != nil {
		log.Printf("error: %v", err)
	}

	log.Printf("res: %v", lock)

	if lock != nil {
		_ = lock.Unlock()
	}

	lock, err = store.Get(context.Background(), "test")
	if err != nil {
		log.Printf("error: %v", err)
	}

	log.Printf("res: %v", lock)

	sigChan := make(chan struct{})

	err = store.Clear(context.Background())
	if err != nil {
		log.Printf("error: %v", err)
	}

	for {
		select {
		case <-sigChan:
			return
		default:
			time.Sleep(time.Second)
		}
	}
}
