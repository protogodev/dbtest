package store

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/RussellLuo/dbtest/builtin"
	sqldb "github.com/RussellLuo/dbtest/builtin/sql"
)

//go:generate dbtest ./store.go Store

type User struct {
	ID    int `dbtest:"-"`
	Name  string
	Sex   string
	Age   int
	Birth time.Time
}

type Store interface {
	CreateUser(user *User) (err error)
	GetUser(name string) (user *User, err error)
	UpdateUser(name string, user *User) (err error)
	DeleteUser(name string) (err error)
}

type DBStore struct {
	db *sql.DB
}

func NewDBStore(dsn string) (*DBStore, error) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	// Always parse time values to time.Time and use UTC.
	cfg.ParseTime = true
	cfg.Loc = time.UTC
	if cfg.Params == nil {
		cfg.Params = map[string]string{"time_zone": "'+00:00'"}
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	return &DBStore{db: db}, nil
}

func (s *DBStore) prepareExec(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *DBStore) CreateUser(user *User) error {
	_, err := s.prepareExec(
		"INSERT INTO user (name, sex, age, birth) VALUES (?, ?, ?, ?)",
		user.Name, user.Sex, user.Age, user.Birth,
	)
	return err
}

func (s *DBStore) GetUser(name string) (*User, error) {
	rows, err := s.db.Query("SELECT * FROM user WHERE name=?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user *User
	for rows.Next() {
		user = new(User)
		if err := rows.Scan(&user.ID, &user.Name, &user.Sex, &user.Age, &user.Birth); err != nil {
			return nil, err
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *DBStore) UpdateUser(name string, user *User) error {
	_, err := s.prepareExec(
		"UPDATE user SET sex=?, age=?, birth=? WHERE name=?",
		user.Sex, user.Age, user.Birth, name,
	)
	return err
}

func (s *DBStore) DeleteUser(name string) error {
	_, err := s.prepareExec(
		"DELETE FROM user WHERE name=?",
		name,
	)
	return err
}

// NewTestee creates a testee for use in tests that are not in this package.
func NewTestee(dsn string) (*builtin.Testee, error) {
	store, err := NewDBStore(dsn)
	if err != nil {
		return nil, err
	}

	return &builtin.Testee{
		SUT: store,
		DB:  sqldb.New(store.db),
	}, nil
}
