
# golock: A Distributed Locking Library

`golock` is a comprehensive Go library designed to facilitate distributed locking across various database systems including Redis, PostgreSQL, MongoDB, MySQL, and Microsoft SQL Server. This library provides a flexible yet robust interface for managing locks in a distributed system, ensuring that only one process can perform certain tasks at any given time.

## Interfaces Overview

The `golock` library defines several interfaces to abstract the locking mechanism across different storage backends:

### `Lock`
- **Unlock() error**: Releases the lock. This method should be called to ensure resources are properly released when the lock is no longer needed.

### `Locker`
- **LockCtx(ctx context.Context, opts ...LockOption) (Lock, error)**: Acquires a lock with the given options within the context's scope. This method returns a `Lock` instance if successful.
- **Lock(opts ...LockOption) (Lock, error)**: Acquires a lock based on the provided options. This is a simplified version without a context.


## Usage Example

Below is a simple usage example of how to acquire and release a lock using the `golock` library with a Redis backend. This example assumes that you have set up the necessary Redis connection and lock store initialization.

```go
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

```

## Contribution

Contributions to `golock` are welcome. Please feel free to submit pull requests, or file issues for bugs, feature requests, or documentation improvements.

## License

`golock` is released under the MIT License. See the LICENSE file in the source directory for more information.
