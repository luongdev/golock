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

func (r *redisLocker) LockCtx(ctx context.Context, opts ...golock.LockOption) (golock.Lock, error) {
	rOpts := make([]lockOption, 0, len(opts))
	for _, opt := range opts {
		if o, ok := opt.(lockOption); ok {
			rOpts = append(rOpts, o)
		}
	}
	lock, err := NewRedisLock(r.store, rOpts...)
	if err != nil {
		return nil, err
	}

	err = r.store.New(ctx, lock)
	if err != nil {
		return nil, err
	}

	return lock, nil
}

func (r *redisLocker) Lock(opts ...golock.LockOption) (golock.Lock, error) {
	return r.LockCtx(context.Background(), opts...)
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
