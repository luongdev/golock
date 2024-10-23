package rdb

import (
	"context"
	"github.com/luongdev/golock"
	internal "github.com/luongdev/golock/internal"
)

type rdbLocker struct {
	*internal.SimpleLocker
}

func (r *rdbLocker) LockCtx(ctx context.Context, opts ...golock.LockOption) (golock.Lock, error) {
	opts = append(opts, internal.WithContext(ctx), internal.WithLockStore(r.LockStore()))
	lock, err := NewRdbLock(opts...)
	if err != nil {
		return nil, err
	}

	return r.SimpleLocker.LockCtx(ctx, lock)
}

func (r *rdbLocker) Lock(opts ...golock.LockOption) (golock.Lock, error) {
	return r.LockCtx(r.Context(), opts...)
}

func NewRdbLocker(store golock.LockStore) golock.Locker {
	return &rdbLocker{
		SimpleLocker: internal.NewSimpleLockerCtx(context.Background(), store),
	}
}

var _ golock.Locker = (*rdbLocker)(nil)
