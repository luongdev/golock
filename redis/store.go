package redis

import (
	"context"
	"fmt"
	"github.com/luongdev/golock"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisLockStore struct {
	rdb *redis.Client

	lockKey string
}

func (r *redisLockStore) New(ctx context.Context, lock golock.Lock) error {
	if rLock, ok := lock.(*redisLock); ok {
		tCtx, cancel := context.WithTimeout(ctx, rLock.lockAtLeast)
		defer cancel()

		setCmd := r.rdb.SetNX(tCtx, r.lockName(rLock.name), lock, rLock.lockAtMost)
		if setCmd.Err() != nil {
			return setCmd.Err()
		}

		if !setCmd.Val() {
			return golock.NewErrLockAlreadyExists(rLock.name)
		}

		return nil
	}

	return golock.NewErrUnsupportedLockType("redis")
}

func (r *redisLockStore) Get(ctx context.Context, name string) (golock.Lock, error) {
	tCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	cmd := r.rdb.TTL(tCtx, r.lockName(name))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	if cmd.Val() == -2 {
		return nil, golock.NewErrLockNotFound(name)
	}

	ttl := cmd.Val()
	if ttl == -1 {
		r.rdb.Expire(tCtx, r.lockName(name), 30*time.Second)
		ttl = 30 * time.Second
	}

	return NewRedisLock(r, WithName(name), WithLockAtMost(ttl))
}

func (r *redisLockStore) Del(ctx context.Context, name string) error {
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cmd := r.rdb.HDel(tCtx, r.lockKey, name)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (r *redisLockStore) Clear(ctx context.Context) error {
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cmd := r.rdb.Del(tCtx, r.lockKey)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (r *redisLockStore) lockName(name string) string {
	return fmt.Sprintf("%s:%s", r.lockKey, name)
}

func NewRedisLockStore(lockKey string) (golock.LockStore, error) {
	s := &redisLockStore{lockKey: lockKey}

	s.rdb = redis.NewClient(&redis.Options{
		Addr:     "34.124.240.249:6979",
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := s.rdb.Ping(ctx)
	if cmd.Err() != nil || cmd.Val() != "PONG" {
		return nil, cmd.Err()
	}

	return s, nil
}

var _ golock.LockStore = (*redisLockStore)(nil)
