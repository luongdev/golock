package main

import (
	"github.com/luongdev/golock/redis"
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
				panic(err)
			}

			defer func() {
				_ = lock.Unlock()
			}()

		}()
	}

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
