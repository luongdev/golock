package rdb

import (
	"github.com/luongdev/golock"
	internal "github.com/luongdev/golock/internal"
	"time"
)

type mongoDoc struct {
	Name        string        `bson:"_id"`
	LockBy      string        `bson:"lock_by"`
	LockTime    time.Time     `bson:"lock_time"`
	LockUntil   time.Time     `bson:"lock_until"`
	LockAtLeast time.Duration `bson:"lock_at_least"`
	LockAtMost  time.Duration `bson:"lock_at_most"`
}

type mongoLock struct {
	*internal.SimpleLock
}

func (r *mongoLock) Doc() *mongoDoc {
	return &mongoDoc{
		Name:        r.Name(),
		LockBy:      r.LockBy(),
		LockTime:    r.LockTime(),
		LockUntil:   r.LockUntil(),
		LockAtLeast: r.LockAtLeast(),
		LockAtMost:  r.LockAtMost(),
	}
}

func NewMongoLock(opts ...golock.LockOption) (golock.Lock, error) {
	sLock, err := internal.NewSimpleLock(opts...)
	if err != nil {
		return nil, err
	}

	return &mongoLock{SimpleLock: sLock}, nil
}

func (r *mongoLock) LockUntil() time.Time {
	return r.LockTime().Add(r.LockAtMost())
}
