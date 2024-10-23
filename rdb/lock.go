package rdb

import (
	"github.com/luongdev/golock"
	internal "github.com/luongdev/golock/internal"
	"time"
)

type rdbEntity struct {
	Name        string `db:"name"`
	LockBy      string `db:"lock_by"`
	LockTime    int64  `db:"lock_time"`
	LockUntil   int64  `db:"lock_until"`
	LockAtLeast string `db:"lock_at_least"`
	LockAtMost  string `db:"lock_at_most"`
}

type rdbLock struct {
	*internal.SimpleLock
}

func NewRdbLock(opts ...golock.LockOption) (golock.Lock, error) {
	sLock, err := internal.NewSimpleLock(opts...)
	if err != nil {
		return nil, err
	}

	return &rdbLock{SimpleLock: sLock}, nil
}

func (r *rdbLock) LockUntil() time.Time {
	return r.LockTime().Add(r.LockAtMost())
}
