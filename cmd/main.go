package main

import (
	"errors"
	"github.com/luongdev/golock"
	"github.com/luongdev/golock/mongo"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

func main() {

	c, err := mongodriver.Connect(options.Client().ApplyURI("mongodb://root:123456@localhost:27017?authSource=admin"))
	if err != nil {
		panic(err)
	}

	store, err := mongo.NewMongoLockStore(c, "distributed_lock")
	if err != nil {
		panic(err)
	}

	locker := mongo.NewMongoLocker(store)
	lock, err := locker.Lock(
		mongo.WithName("lock_name"),
		mongo.WithLockBy("process01"),
		mongo.WithLockAtLeast(time.Second*5),
		mongo.WithLockAtMost(time.Second*10),
	)

	if err != nil {
		if errors.Is(err, golock.ErrLockAlreadyExists) {
			panic("lock already exists")
		}
	}

	defer lock.Unlock()

	select {}
}
