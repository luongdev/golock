package mongo

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/luongdev/golock"
	_ "github.com/microsoft/go-mssqldb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
	"time"
)

type mongoLockStore struct {
	c  *mongo.Client
	db *mongo.Database
}

func (s *mongoLockStore) New(ctx context.Context, lock golock.Lock) error {
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if mLock, ok := lock.(*mongoLock); ok {
		col := s.db.Collection("Locker")
		one, err := col.InsertOne(tCtx, mLock.Doc())
		if err != nil {
			return golock.NewErrLockAlreadyExists(mLock.Name())
		}
		log.Printf("Inserted a single document: %v", one.InsertedID)
		return nil
	}
	return golock.NewErrUnsupportedLockType("rdb")
}

func (s *mongoLockStore) Get(ctx context.Context, name string) (golock.Lock, error) {
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	col := s.db.Collection("Locker")
	r := col.FindOne(tCtx, bson.M{"_id": name})
	if r.Err() != nil {
		return nil, golock.NewErrLockNotFound(name)
	}

	var doc mongoDoc
	err := r.Decode(&doc)
	if err != nil {
		return nil, golock.NewErrLockNotFound(name)
	}

	return NewMongoLock(
		WithName(doc.Name),
		WithLockAtLeast(doc.LockAtLeast),
		WithLockAtMost(doc.LockAtMost),
		WithLockBy(doc.LockBy),
		WithLockTime(doc.LockTime),
		WithContext(ctx),
		WithLockStore(s),
	)
}

func (s *mongoLockStore) Del(ctx context.Context, name string) error {

	return nil
}

func (s *mongoLockStore) Clear(ctx context.Context) error {

	return nil
}

func NewMongoLockStore(c *mongo.Client, dbName string) (golock.LockStore, error) {
	db := c.Database(dbName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.CreateCollection(ctx, "Locker")
	if err != nil {
		log.Printf("warn: %v", err)
	}

	return &mongoLockStore{c: c, db: db}, nil
}

var _ golock.LockStore = (*mongoLockStore)(nil)
