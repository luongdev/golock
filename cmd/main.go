package main

import (
	"context"
	libsqlx "github.com/jmoiron/sqlx"
	"github.com/luongdev/golock/rdb"
	"github.com/luongdev/golock/redis"
	"log"
	"time"
)

func main() {

	db, err := libsqlx.Connect("postgres", "user=postgres password=Default#Postgres@6699 dbname=freeswitch host=192.168.13.137 port=15432 sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	store, err := rdb.NewSqlxLockStore(db)
	if err != nil {
		log.Fatalln(err)
	}

	//store, err := redis.NewRedisLockStore("DefaultLockKey", libredis.NewClient(&libredis.Options{}))
	//if err != nil {
	//	panic(err)
	//}
	//

	locker := rdb.NewRdbLocker(store)

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

	//err = store.Clear(context.Background())
	//if err != nil {
	//	log.Printf("error: %v", err)
	//}

	for {
		select {
		case <-sigChan:
			return
		default:
			time.Sleep(time.Second)
		}
	}
}
