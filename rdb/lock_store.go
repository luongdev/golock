package rdb

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/luongdev/golock"
	_ "github.com/microsoft/go-mssqldb"
	"log"
	"time"
)

var schema = `
CREATE TABLE Locker (
    name 			VARCHAR(255) PRIMARY KEY,
    lock_by 		VARCHAR(255),
    lock_time 		bigint,
    lock_until 		bigint,
    lock_at_most 	varchar(16),
    lock_at_least 	varchar(16)
);
`

type rdbLockStore struct {
	db *sqlx.DB
}

func (s *rdbLockStore) New(ctx context.Context, lock golock.Lock) error {
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if rLock, ok := lock.(*rdbLock); ok {
		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(tCtx,
			"INSERT INTO Locker (name, lock_by, lock_time, lock_until, lock_at_most, lock_at_least) VALUES ($1, $2, $3, $4, $5, $6);",
			rLock.Name(), rLock.LockBy(), rLock.LockTime().Unix(), rLock.LockUntil().Unix(), rLock.LockAtMost().String(), rLock.LockAtLeast().String(),
		)

		if err != nil {
			return golock.NewErrLockAlreadyExists(rLock.Name())
		}

		err = tx.Commit()
		if err != nil {
			return tx.Rollback()
		}

		return nil
	}

	return golock.NewErrUnsupportedLockType("rdb")
}

func (s *rdbLockStore) Get(ctx context.Context, name string) (golock.Lock, error) {
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var rLocks []rdbEntity
	err := s.db.SelectContext(tCtx, &rLocks, "SELECT * FROM Locker WHERE name = $1;", name)
	if err != nil {
		return nil, err
	}

	if errors.Is(err, sql.ErrNoRows) || len(rLocks) == 0 {
		return nil, golock.NewErrLockNotFound(name)
	}

	rLock := rLocks[0]

	lockAtLeast, err := time.ParseDuration(rLock.LockAtLeast)
	if err != nil {
		lockAtLeast = time.Second * 2
	}
	lockAtMost, err := time.ParseDuration(rLock.LockAtMost)
	if err != nil {
		lockAtMost = time.Second * 30
	}

	return NewRdbLock(
		WithName(rLock.Name),
		WithLockTime(time.Unix(rLock.LockTime, 0)),
		WithLockAtLeast(lockAtLeast),
		WithLockAtMost(lockAtMost),
		WithLockStore(s),
		WithContext(ctx),
	)
}

func (s *rdbLockStore) Del(ctx context.Context, name string) error {
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(tCtx, "DELETE FROM Locker WHERE name = $1;", name)
	return err
}

func (s *rdbLockStore) Clear(ctx context.Context) error {
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(tCtx, "DELETE FROM Locker;")
	return err
}

func NewRdbLockStore(db *sqlx.DB) (golock.LockStore, error) {
	_, err := db.Exec(schema)
	if err != nil {
		log.Printf("warn: %v", err)
	}

	r, err := db.Exec("DELETE FROM Locker WHERE lock_until <= $1;", time.Now().Unix())
	if err != nil {
		log.Printf("warn: %v", err)
	} else {
		if n, err := r.RowsAffected(); err == nil && n > 0 {
			log.Printf("info: %d expired locks removed", n)
		}
	}

	return &rdbLockStore{db: db}, nil
}

var _ golock.LockStore = (*rdbLockStore)(nil)
