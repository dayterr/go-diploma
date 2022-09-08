package storage

import (
	"context"
	"database/sql"
	"github.com/dayterr/go-diploma/cmd/gophermart/handlers"
	"log"
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

func (us UserStorage) AddUser(user handlers.UserModel) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := us.DB.ExecContext(ctx,
		`INSERT INTO users (login, password) VALUES ($1, $2) RETUNRNING id`,
		user.Name, user.Password)
	if err != nil {
		log.Println("user creation error:", err)
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, err
}
