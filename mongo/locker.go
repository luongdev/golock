package mongo

import (
	"context"
	"github.com/luongdev/golock"
	internal "github.com/luongdev/golock/internal"
)

type mongoLocker struct {
	*internal.SimpleLocker
}

func (r *mongoLocker) LockCtx(ctx context.Context, opts ...golock.LockOption) (golock.Lock, error) {
	opts = append(opts, internal.WithContext(ctx), internal.WithLockStore(r.LockStore()))
	lock, err := NewMongoLock(opts...)
	if err != nil {
		return nil, err
	}

	return r.SimpleLocker.LockCtx(ctx, lock)
}

func (r *mongoLocker) Lock(opts ...golock.LockOption) (golock.Lock, error) {
	return r.LockCtx(r.Context(), opts...)
}

func NewMongoLocker(store golock.LockStore) golock.Locker {
	return &mongoLocker{
		SimpleLocker: internal.NewSimpleLockerCtx(context.Background(), store),
	}
}

var _ golock.Locker = (*mongoLocker)(nil)
