package storage

import (
	"database/sql"
)

type Storager interface {
	AddUser(user UserModel) (int64, error)
	GetUser(user UserModel) (int64, error)
}

type UserStorage struct {
	DB *sql.DB
	DSN string
}

type UserModel struct {
	ID int
	Name string
	Password string
}