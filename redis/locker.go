package redis

import (
	"context"
	"github.com/luongdev/golock"
	internal "github.com/luongdev/golock/internal"
)

type redisLocker struct {
	*internal.SimpleLocker
}

func (r *redisLocker) LockCtx(ctx context.Context, opts ...golock.LockOption) (golock.Lock, error) {
	opts = append(opts, internal.WithContext(ctx), internal.WithLockStore(r.LockStore()))
	lock, err := NewRedisLock(opts...)
	if err != nil {
		return nil, err
	}

	return r.SimpleLocker.LockCtx(ctx, lock)
}

func (r *redisLocker) Lock(opts ...golock.LockOption) (golock.Lock, error) {
	return r.LockCtx(r.Context(), opts...)
}

func NewRedisLocker(store golock.LockStore) golock.Locker {
	return &redisLocker{
		SimpleLocker: internal.NewSimpleLockerCtx(context.Background(), store),
	}
}

var _ golock.Locker = (*redisLocker)(nil)
