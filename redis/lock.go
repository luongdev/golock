package redis

import (
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/luongdev/golock"
	"time"
)

type lockOption func(*redisLock) error

func (l lockOption) Apply(lock golock.Lock) error {
	if rLock, ok := lock.(*redisLock); ok {
		return l(rLock)
	}
	return nil
}

var _ golock.LockOption = (lockOption)(nil)

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

func (r *redisLock) String() string {
	return fmt.Sprintf(
		"redisLock{name=%s, lockTime=%v, lockAtLeast=%v, lockAtMost=%v}",
		r.name, r.lockTime, r.lockAtLeast, r.lockAtMost,
	)
}

func (r *redisLock) MarshalBinary() (data []byte, err error) {
	s := jsonLock{Name: r.name, LockTime: r.lockTime, LockAtLeast: r.lockAtLeast, LockAtMost: r.lockAtMost}
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

func NewRedisLock(store *redisLockStore, opts ...lockOption) (golock.Lock, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	l := &redisLock{
		name:        uid.String(),
		lockTime:    time.Now(),
		lockAtLeast: time.Second * 2,
		lockAtMost:  time.Second * 30,
		store:       store,
	}

	for _, opt := range opts {
		if err = opt(l); err != nil {
			return nil, err
		}
	}

	if l.ctx == nil {
		l.ctx, l.cancel = context.WithTimeout(context.Background(), l.lockAtMost)
	}

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

func WithName(name string) golock.LockOption {
	return lockOption(func(l *redisLock) error {
		l.name = name
		return nil
	})
}

func WithContext(ctx context.Context) golock.LockOption {
	return lockOption(func(l *redisLock) error {
		l.ctx, l.cancel = context.WithCancel(ctx)
		return nil
	})
}

func WithLockTime(lockTime time.Time) golock.LockOption {
	return lockOption(func(l *redisLock) error {
		l.lockTime = lockTime
		return nil
	})
}

func WithLockAtLeast(lockAtLeast time.Duration) golock.LockOption {
	return lockOption(func(l *redisLock) error {
		l.lockAtLeast = lockAtLeast
		return nil
	})
}

func WithLockAtMost(lockAtMost time.Duration) golock.LockOption {
	return lockOption(func(l *redisLock) error {
		l.lockAtMost = lockAtMost
		return nil
	})
}

var _ golock.Lock = (*redisLock)(nil)
var _ encoding.BinaryUnmarshaler = (*redisLock)(nil)
var _ encoding.BinaryMarshaler = (*redisLock)(nil)
var _ fmt.Stringer = (*redisLock)(nil)
