package golock

import "errors"

var (
	ErrLockAlreadyExists   error = &errLockAlreadyExists{}
	ErrUnsupportedLockType error = &unsupportedLockType{}
	ErrLockNotFound        error = &errLockNotFound{}
)

type errLockAlreadyExists struct {
	name string
}

func (e *errLockAlreadyExists) Is(target error) bool {
	var ee *errLockAlreadyExists
	return errors.As(target, &ee)
}

func (e *errLockAlreadyExists) Error() string {
	return "lock already exists: " + e.name
}

func NewErrLockAlreadyExists(name string) error {
	return &errLockAlreadyExists{name: name}
}

var _ error = (*errLockAlreadyExists)(nil)

type unsupportedLockType struct {
	lockType string
}

func (e *unsupportedLockType) Is(target error) bool {
	var ee *unsupportedLockType
	return errors.As(target, &ee)
}

func (e *unsupportedLockType) Error() string {
	return "supported only lock type: " + e.lockType
}

func NewErrUnsupportedLockType(lockType string) error {
	return &unsupportedLockType{lockType: lockType}
}

type errLockNotFound struct {
	name string
}

func (e *errLockNotFound) Is(target error) bool {
	var ee *errLockNotFound
	return errors.As(target, &ee)
}

func (e *errLockNotFound) Error() string {
	return "lock not found: " + e.name
}

func NewErrLockNotFound(name string) error {
	return &errLockNotFound{name: name}
}
