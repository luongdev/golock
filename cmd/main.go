package main

import (
	"context"
	rdb "github.com/luongdev/golock/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"time"
)

func main() {

	c, err := mongo.Connect(options.Client().ApplyURI("mongodb://root:Abcd54321@192.168.13.54:27017/test?authSource=admin"))
	if err != nil {
		panic(err)
	}

	store, err := rdb.NewMongoLockStore(c, "test")
	if err != nil {
		panic(err)
	}

	//log.Printf("MongoDB lock store created %v\n", store)
	//
	//mLock, err := rdb.NewMongoLock(rdb.WithName("test"), rdb.WithLockStore(store))
	//if err != nil {
	//	panic(err)
	//}
	//for i := 0; i < 3; i++ {
	//	go func() {
	//		err = store.New(context.Background(), mLock)
	//		if err != nil {
	//			return
	//		}
	//	}()
	//}

	mLock, err := store.Get(context.Background(), "test")
	if err != nil {
		panic(err)
	}

	log.Printf("MongoDB lock created %v\n", mLock)

	sigChan := make(chan struct{})
	for {
		select {
		case <-sigChan:
			return
		default:
			time.Sleep(time.Second)
		}
	}
}
