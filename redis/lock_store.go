package redis

import (
	"context"
	"encoding/json"
	"github.com/luongdev/golock"
	internal "github.com/luongdev/golock/internal"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisLockStore struct {
	*internal.SimpleLockStore

	rdb *redis.Client
}

func (r *redisLockStore) New(ctx context.Context, lock golock.Lock) error {
	if rLock, ok := lock.(*redisLock); ok {
		tCtx, cancel := context.WithTimeout(ctx, rLock.LockAtLeast())
		defer cancel()

		setCmd := r.rdb.SetNX(tCtx, r.LockName(rLock.Name()), lock, rLock.LockAtMost())
		if setCmd.Err() != nil {
			return setCmd.Err()
		}

		if !setCmd.Val() {
			return golock.NewErrLockAlreadyExists(rLock.Name())
		}

		return nil
	}

	return golock.NewErrUnsupportedLockType("redis")
}

func (r *redisLockStore) Get(ctx context.Context, name string) (golock.Lock, error) {
	tCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	rLock := &redisLock{}
	cmd := r.rdb.Get(tCtx, r.LockName(name))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	if b, err := cmd.Bytes(); err != nil {
		return nil, golock.NewErrLockNotFound(name)
	} else {
		var j redisJson
		if err = json.Unmarshal(b, &j); err != nil {
			return nil, err
		}

		sLock, err := internal.NewSimpleLock(
			WithName(j.Name),
			WithLockTime(j.LockTime),
			WithLockAtLeast(j.LockAtLeast),
			WithLockAtMost(j.LockAtMost),
			WithLockStore(r),
			WithContext(ctx),
		)
		if err != nil {
			return nil, err
		}

		rLock.SimpleLock = sLock

		if time.Now().After(rLock.LockTime().Add(rLock.LockAtMost())) {
			_ = r.Del(ctx, name)
			return nil, golock.NewErrLockNotFound(name)
		}

		return rLock, nil
	}
}

func (r *redisLockStore) Del(ctx context.Context, name string) error {
	tCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	cmd := r.rdb.ExpireAt(tCtx, r.LockName(name), time.Now())
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (r *redisLockStore) Clear(ctx context.Context) error {
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cmd := r.rdb.Keys(tCtx, r.LockName("*"))
	if cmd.Err() != nil {
		return cmd.Err()
	}

	for _, key := range cmd.Val() {
		delCmd := r.rdb.Del(tCtx, key)
		if delCmd.Err() != nil {
			return delCmd.Err()
		}
	}

	return nil
}

func NewRedisLockStore(lockKey string, rdb *redis.Client) (golock.LockStore, error) {
	s := &redisLockStore{
		rdb:             rdb,
		SimpleLockStore: internal.NewSimpleLockStore(lockKey),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := s.rdb.Ping(ctx)
	if cmd.Err() != nil || cmd.Val() != "PONG" {
		return nil, cmd.Err()
	}

	return s, nil
}

var _ golock.LockStore = (*redisLockStore)(nil)
