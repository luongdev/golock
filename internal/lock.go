package golock

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/luongdev/golock"
	"os"
	"time"
)

type SimpleLockOption func(*SimpleLock) error

func (s SimpleLockOption) Apply(lock golock.Lock) error {
	if sLock, ok := lock.(*SimpleLock); ok {
		return s(sLock)
	}
	return nil
}

var _ golock.LockOption = (SimpleLockOption)(nil)

type SimpleLock struct {
	name        string
	lockBy      string
	lockTime    time.Time
	lockAtLeast time.Duration
	lockAtMost  time.Duration

	ctx    context.Context
	cancel context.CancelFunc
	store  golock.LockStore
}

func (s *SimpleLock) Name() string {
	return s.name
}

func (s *SimpleLock) LockBy() string {
	return s.lockBy
}

func (s *SimpleLock) LockTime() time.Time {
	return s.lockTime
}

func (s *SimpleLock) LockAtLeast() time.Duration {
	return s.lockAtLeast
}

func (s *SimpleLock) LockAtMost() time.Duration {
	return s.lockAtMost
}

func (s *SimpleLock) LockStore() golock.LockStore {
	return s.store
}

func (s *SimpleLock) Context() (context.Context, context.CancelFunc) {
	return s.ctx, s.cancel
}

func NewSimpleLock(opts ...golock.LockOption) (*SimpleLock, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	l := &SimpleLock{
		name:        uid.String(),
		lockTime:    time.Now(),
		lockAtLeast: 2 * time.Second,
		lockAtMost:  15 * time.Second,
		lockBy:      hostName,
	}

	for _, opt := range opts {
		if err = opt.Apply(l); err != nil {
			return nil, err
		}
	}

	if l.ctx == nil {
		l.ctx, l.cancel = context.WithTimeout(context.Background(), l.lockAtMost)
	}

	if l.store == nil {
		return nil, fmt.Errorf("lock store is required")
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

func (s *SimpleLock) String() string {
	return fmt.Sprintf("SimpleLock{"+
		"name=%v, "+
		"lockTime=%v, "+
		"lockAtLeast=%v, "+
		"lockAtMost=%v}",
		s.name, s.lockTime, s.lockAtLeast, s.lockAtMost)
}

func (s *SimpleLock) Unlock() error {
	return s.store.Del(context.Background(), s.name)
}

var _ golock.Lock = (*SimpleLock)(nil)

func WithName(name string) golock.LockOption {
	return SimpleLockOption(func(l *SimpleLock) error {
		l.name = name
		return nil
	})
}

func WithLockBy(lockBy string) golock.LockOption {
	return SimpleLockOption(func(l *SimpleLock) error {
		l.lockBy = lockBy
		return nil
	})
}

func WithContext(ctx context.Context) golock.LockOption {
	return SimpleLockOption(func(l *SimpleLock) error {
		l.ctx, l.cancel = context.WithCancel(ctx)
		return nil
	})
}

func WithLockTime(lockTime time.Time) golock.LockOption {
	return SimpleLockOption(func(l *SimpleLock) error {
		l.lockTime = lockTime
		return nil
	})
}

func WithLockAtLeast(lockAtLeast time.Duration) golock.LockOption {
	return SimpleLockOption(func(l *SimpleLock) error {
		l.lockAtLeast = lockAtLeast
		return nil
	})
}

func WithLockAtMost(lockAtMost time.Duration) golock.LockOption {
	return SimpleLockOption(func(l *SimpleLock) error {
		l.lockAtMost = lockAtMost
		return nil
	})
}

func WithLockStore(store golock.LockStore) golock.LockOption {
	return SimpleLockOption(func(l *SimpleLock) error {
		l.store = store
		return nil
	})
}
