package redis

import (
	"context"
	"github.com/luongdev/golock"
)

type redisLocker struct {
	store *redisLockStore

	ctx    context.Context
	cancel context.CancelFunc
}

func (r *redisLocker) LockCtx(ctx context.Context, name string) (golock.Lock, error) {
	lock, err := NewRedisLock(r.store, WithName(name))
	if err != nil {
		return nil, err
	}

	err = r.store.New(ctx, lock)
	if err != nil {
		return nil, err
	}

	return lock, nil
}

func (r *redisLocker) Lock(name string) (golock.Lock, error) {
	return r.LockCtx(context.Background(), name)
}

func NewRedisLocker(store golock.LockStore) golock.Locker {
	ctx, cancel := context.WithCancel(context.Background())
	return &redisLocker{
		store:  store.(*redisLockStore),
		ctx:    ctx,
		cancel: cancel,
	}
}

var _ golock.Locker = (*redisLocker)(nil)
