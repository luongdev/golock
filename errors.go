package golock

type ErrLockAlreadyExists struct {
	name string
}

func (e *ErrLockAlreadyExists) Error() string {
	return "lock already exists: " + e.name
}

func NewErrLockAlreadyExists(name string) *ErrLockAlreadyExists {
	return &ErrLockAlreadyExists{name: name}
}

var _ error = (*ErrLockAlreadyExists)(nil)

type UnsupportedLockType struct {
	lockType string
}

func (e *UnsupportedLockType) Error() string {
	return "supported only lock type: " + e.lockType
}

func NewUnsupportedLockType(lockType string) *UnsupportedLockType {
	return &UnsupportedLockType{lockType: lockType}
}
