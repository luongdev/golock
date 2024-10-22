package golock

import (
	"fmt"
)

type SimpleLockStore struct {
	lockKey string
}

func (s *SimpleLockStore) LockKey() string {
	return s.lockKey
}

func (s *SimpleLockStore) LockName(name string) string {
	return fmt.Sprintf("%s:%s", s.lockKey, name)
}

func NewSimpleLockStore(lockKey string) *SimpleLockStore {
	return &SimpleLockStore{lockKey: lockKey}
}
