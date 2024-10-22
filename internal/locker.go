package golock

import (
	"context"
	"github.com/luongdev/golock"
)

type SimpleLocker struct {
	store golock.LockStore

	ctx context.Context
}

func (r *SimpleLocker) Context() context.Context {
	return r.ctx
}

func (r *SimpleLocker) LockStore() golock.LockStore {
	return r.store
}

func (r *SimpleLocker) LockCtx(ctx context.Context, lock golock.Lock) (golock.Lock, error) {
	err := r.store.New(ctx, lock)
	if err != nil {
		return nil, err
	}

	return lock, nil
}

func (r *SimpleLocker) Lock(lock golock.Lock) (golock.Lock, error) {
	return r.LockCtx(r.Context(), lock)
}

func NewSimpleLockerCtx(ctx context.Context, store golock.LockStore) *SimpleLocker {
	locker := &SimpleLocker{store: store, ctx: ctx}

	return locker
}

func NewSimpleLocker(store golock.LockStore) *SimpleLocker {
	return NewSimpleLockerCtx(context.Background(), store)
}
