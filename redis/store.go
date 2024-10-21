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

	l, err := NewRedisLock(r)
	if err != nil {
		return nil, err
	}

	rLock := l.(*redisLock)

	cmd := r.rdb.Get(tCtx, r.lockName(name))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	if b, err := cmd.Bytes(); err != nil {
		return nil, golock.NewErrLockNotFound(name)
	} else if err := rLock.UnmarshalBinary(b); err != nil {
		return nil, golock.NewErrLockNotFound(name)
	}

	if time.Now().After(rLock.lockTime.Add(rLock.lockAtMost)) {
		_ = r.Del(ctx, name)
		return nil, golock.NewErrLockNotFound(name)
	}
	rLock.ctx, rLock.cancel = context.WithTimeout(ctx, rLock.lockAtMost)

	return l, nil
}

func (r *redisLockStore) Del(ctx context.Context, name string) error {
	tCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	cmd := r.rdb.ExpireAt(tCtx, r.lockName(name), time.Now())
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (r *redisLockStore) Clear(ctx context.Context) error {
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cmd := r.rdb.Keys(tCtx, r.lockName("*"))
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

func (r *redisLockStore) lockName(name string) string {
	return fmt.Sprintf("%s:%s", r.lockKey, name)
}

func NewRedisLockStore(lockKey string, rdb *redis.Client) (golock.LockStore, error) {
	s := &redisLockStore{lockKey: lockKey, rdb: rdb}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := s.rdb.Ping(ctx)
	if cmd.Err() != nil || cmd.Val() != "PONG" {
		return nil, cmd.Err()
	}

	return s, nil
}

var _ golock.LockStore = (*redisLockStore)(nil)
