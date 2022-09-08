package storage

import (
	"context"
	"database/sql"
	"time"
)

func NewUserDB(dsn string) (UserStorage, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return UserStorage{}, err
	}
	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS users (id serial PRIMARY KEY, login text UNIQUE NOT NULL, password text NOT NULL);`)
	if err != nil {
		return UserStorage{}, err
	}

	return UserStorage{
		DB:           db,
		DSN:          dsn,
	}, nil
}
