package redis

import (
	"context"
	"github.com/luongdev/golock"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisLockStore struct {
	rdb *redis.Client

	lockKey   string
	lockValue bool
}

func (r *redisLockStore) New(lock golock.Lock) error {
	if rLock, ok := lock.(*redisLock); ok {
		ctx, cancel := context.WithTimeout(context.Background(), rLock.lockAtLeast)
		defer cancel()

		setCmd := r.rdb.HSetNX(ctx, r.lockKey, rLock.name, r.lockValue)
		if setCmd.Err() != nil {
			return setCmd.Err()
		}

		if !setCmd.Val() {
			return golock.NewErrLockAlreadyExists(rLock.name)
		}

		expCmd := r.rdb.HExpire(ctx, r.lockKey, rLock.lockAtMost, rLock.name)
		if expCmd.Err() != nil {
			_ = r.rdb.HDel(ctx, r.lockKey, rLock.name)
			return expCmd.Err()
		}

		return nil
	}

	return golock.NewUnsupportedLockType("redis")
}

func (r *redisLockStore) Get(name string) (golock.Lock, error) {
	//TODO implement me
	panic("implement me")
}

func (r *redisLockStore) Del(name string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := r.rdb.HDel(ctx, r.lockKey, name)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (r *redisLockStore) Clear() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := r.rdb.Del(ctx, r.lockKey)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func NewRedisLockStore(lockKey string) (golock.LockStore, error) {
	s := &redisLockStore{lockKey: lockKey, lockValue: true}

	s.rdb = redis.NewClient(&redis.Options{
		Addr:     "34.124.240.249:6979",
		Password: "",
		DB:       0,
	})

	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	//cmd := s.rdb.HSetNX(ctx, s.lockKey, "init", true)
	//if cmd.Err() != nil {
	//	return nil, cmd.Err()
	//}

	return s, nil
}

var _ golock.LockStore = (*redisLockStore)(nil)
