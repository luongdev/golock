package mongo

import (
	"context"
	"github.com/luongdev/golock"
	internal "github.com/luongdev/golock/internal"
	"time"
)

func WithName(name string) golock.LockOption {
	return internal.WithName(name)
}

func WithLockAtLeast(d time.Duration) golock.LockOption {
	return internal.WithLockAtLeast(d)
}

func WithLockAtMost(d time.Duration) golock.LockOption {
	return internal.WithLockAtMost(d)
}

func WithLockBy(s string) golock.LockOption {
	return internal.WithLockBy(s)
}

func WithLockTime(t time.Time) golock.LockOption {
	return internal.WithLockTime(t)
}

func WithContext(ctx context.Context) golock.LockOption {
	return internal.WithContext(ctx)
}

func WithLockStore(store golock.LockStore) golock.LockOption {
	return internal.WithLockStore(store)
}
