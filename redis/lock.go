package redis

import (
	"context"
	"encoding"
	"encoding/json"
	"github.com/luongdev/golock"
	"time"
)

type LockOption func(*redisLock) error

type jsonLock struct {
	Name        string
	LockTime    time.Time
	LockAtLeast time.Duration
	LockAtMost  time.Duration
}

type redisLock struct {
	ctx    context.Context
	cancel context.CancelFunc

	name        string
	lockTime    time.Time
	lockAtLeast time.Duration
	lockAtMost  time.Duration

	store *redisLockStore
}

func (r *redisLock) MarshalBinary() (data []byte, err error) {
	s := jsonLock{
		Name:        r.name,
		LockTime:    r.lockTime,
		LockAtLeast: r.lockAtLeast,
	}
	return json.Marshal(s)
}

func (r *redisLock) UnmarshalBinary(data []byte) error {
	var s jsonLock
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	r.name = s.Name
	r.lockTime = s.LockTime
	r.lockAtLeast = s.LockAtLeast
	r.lockAtMost = s.LockAtMost

	return nil
}

func (r *redisLock) Unlock() error {
	return r.store.Del(r.ctx, r.name)
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

	for _, opt := range opts {
		if err := opt(l); err != nil {
			return nil, err
		}
	}

	l.ctx, l.cancel = context.WithTimeout(context.Background(), l.lockAtMost)

	go func() {
		for {
			select {
			case <-l.ctx.Done():
				return
			default:
				if time.Now().After(l.lockTime.Add(l.lockAtMost)) {
					_ = l.store.Del(l.ctx, l.name)
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

func WithLockAtMost(lockAtMost time.Duration) LockOption {
	return func(l *redisLock) error {
		l.lockAtMost = lockAtMost
		return nil
	}
}

var _ golock.Lock = (*redisLock)(nil)
var _ encoding.BinaryUnmarshaler = (*redisLock)(nil)
var _ encoding.BinaryMarshaler = (*redisLock)(nil)
