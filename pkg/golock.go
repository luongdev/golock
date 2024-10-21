package pkg

import (
	"github.com/luongdev/golock"
	"github.com/luongdev/golock/internal"
	"time"
)

type LockOption = internal.LockOption

func WithLockTime(t time.Time) LockOption {
	return internal.WithLockTime(t)
}

func WithLockAtLeast(d time.Duration) LockOption {
	return internal.WithLockAtLeast(d)
}

func WithLockAtMost(d time.Duration) LockOption {
	return internal.WithLockAtMost(d)
}

func WithLockStore(store golock.LockStore) LockOption {
	return internal.WithLockStore(store)
}
