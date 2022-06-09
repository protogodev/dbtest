package builtin

import (
	"testing"

	"github.com/protogodev/dbtest/spec"
)

type DB interface {
	Insert(data map[string]spec.Rows) error
	Delete(keys ...string) error
	Select(query string) (spec.Rows, error)
	Close() error
}

type Fixture struct {
	t    *testing.T
	db   DB
	data map[string]spec.Rows
}

func NewFixture(t *testing.T, db DB, data map[string]spec.Rows) *Fixture {
	return &Fixture{t: t, db: db, data: data}
}

func (f *Fixture) SetUp() {
	if err := f.db.Insert(f.data); err != nil {
		f.t.Fatalf("err: %v", err)
	}
}

func (f *Fixture) TearDown() {
	var keys []string
	for tableName := range f.data {
		keys = append(keys, tableName)
	}

	if err := f.db.Delete(keys...); err != nil {
		f.t.Fatalf("err: %v", err)
	}
}

func (f *Fixture) Reset() {
	f.TearDown()
	f.SetUp()
}

func (f *Fixture) Query(query string) (result spec.Rows) {
	result, err := f.db.Select(query)
	if err != nil {
		f.t.Fatalf("err: %v", err)
	}
	return result
}
