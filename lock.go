package golock

import (
	"context"
)

type Lock interface {
	Unlock() error
}

type Locker interface {
	LockCtx(ctx context.Context, name string) (Lock, error)
	Lock(name string) (Lock, error)
}

type LockStore interface {
	New(lock Lock) error
	Get(name string) (Lock, error)
	Del(name string) error
	Clear() error
}
