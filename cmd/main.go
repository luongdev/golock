package main

import (
	"context"
	"github.com/luongdev/golock/redis"
	"log"
	"time"
)

func main() {
	store, err := redis.NewRedisLockStore("DefaultLockKey")
	if err != nil {
		panic(err)
	}

	locker := redis.NewRedisLocker(store)

	for i := 0; i < 3; i++ {
		go func() {
			lock, err := locker.Lock("test")
			if err != nil {
				log.Printf("error: %v", err)
			}

			defer func() {
				if lock != nil {
					_ = lock.Unlock()
				}
			}()

		}()
	}

	time.Sleep(time.Second * 2)

	res, err := store.Get(context.Background(), "test")
	if err != nil {
		log.Printf("error: %v", err)
	}

	log.Printf("res: %v", res)

	sigChan := make(chan struct{})
	//time.AfterFunc(time.Second*5, func() {
	//	_ = lock.Unlock()
	//
	//	sigChan <- struct{}{}
	//})

	for {
		select {
		case <-sigChan:
			return
		default:
			time.Sleep(time.Second)
		}
	}
	//s.Clear()
}
