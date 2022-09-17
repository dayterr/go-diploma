package storage

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
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
	log.Println("table users created")
	return UserStorage{
		DB:           db,
		DSN:          dsn,
	}, nil
}

func (us UserStorage) AddUser(user UserModel) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("writing user to database")
	res, err := us.DB.ExecContext(ctx,
		`INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id`,
		user.Name, user.Password)
	if err != nil {
		log.Println("user creation error:", err)
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, err
}

func (us UserStorage) GetUser(user UserModel) (int64, error) {
	var userID int64

	log.Println("getting user from database")
	res := us.DB.QueryRow(`SELECT id FROM users WHERE login = $1`, user.Name)
	err := res.Scan(&userID)
	if err != nil {
		log.Println("getting id error", err)
		return 0, err
	}
	return userID, err
}
