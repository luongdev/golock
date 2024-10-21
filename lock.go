package golock

import (
	"context"
)

type LockOption interface {
	Apply(Lock) error
}

type Lock interface {
	Unlock() error
}

type Locker interface {
	LockCtx(ctx context.Context, opts ...LockOption) (Lock, error)
	Lock(opts ...LockOption) (Lock, error)
}

type LockStore interface {
	New(ctx context.Context, lock Lock) error
	Get(ctx context.Context, name string) (Lock, error)
	Del(ctx context.Context, name string) error
	Clear(ctx context.Context) error
}
