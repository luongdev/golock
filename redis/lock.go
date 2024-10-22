package redis

import (
	"encoding"
	"encoding/json"
	"github.com/luongdev/golock"
	internal "github.com/luongdev/golock/internal"
	"time"
)

type redisJson struct {
	Name        string
	LockTime    time.Time
	LockAtLeast time.Duration
	LockAtMost  time.Duration
}

type redisLock struct {
	*internal.SimpleLock
}

func NewRedisLock(opts ...golock.LockOption) (golock.Lock, error) {
	sLock, err := internal.NewSimpleLock(opts...)
	if err != nil {
		return nil, err
	}

	return &redisLock{SimpleLock: sLock}, nil
}

func (s *redisLock) MarshalBinary() (data []byte, err error) {
	return json.Marshal(&redisJson{
		Name:        s.Name(),
		LockTime:    s.LockTime(),
		LockAtLeast: s.LockAtLeast(),
		LockAtMost:  s.LockAtMost(),
	})
}

var _ encoding.BinaryMarshaler = (*redisLock)(nil)
var _ golock.Lock = (*redisLock)(nil)
