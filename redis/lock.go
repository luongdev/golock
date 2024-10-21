package redis

import (
	"context"
	"github.com/luongdev/golock"
	"time"
)

type LockOption func(*redisLock) error

type redisLock struct {
	ctx    context.Context
	cancel context.CancelFunc

	name        string
	lockTime    time.Time
	lockAtLeast time.Duration
	lockAtMost  time.Duration

	store *redisLockStore
}

func (r *redisLock) Unlock() error {
	return r.store.Del(r.name)
}

func (r *redisLock) getLockTimes() (time.Time, time.Time) {
	return r.lockTime, r.lockTime.Add(r.lockAtMost)
}

func NewRedisLock(store *redisLockStore, opts ...LockOption) (golock.Lock, error) {
	l := &redisLock{
		lockTime:    time.Now(),
		lockAtLeast: time.Second * 2,
		lockAtMost:  time.Second * 30,
		store:       store,
	}

	l.ctx, l.cancel = context.WithCancel(context.Background())
	for _, opt := range opts {
		if err := opt(l); err != nil {
			return nil, err
		}
	}

	go func() {
		for {
			select {
			case <-l.ctx.Done():
				return
			default:
				if time.Now().After(l.lockTime.Add(l.lockAtMost)) {
					_ = l.store.Del(l.name)
					return
				}
				time.Sleep(l.lockAtLeast)
			}
		}
	}()

	return l, nil
}

func WithName(name string) LockOption {
	return func(l *redisLock) error {
		l.name = name
		return nil
	}
}

var _ golock.Lock = (*redisLock)(nil)
