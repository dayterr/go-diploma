package storage

import (
	"database/sql"
	"github.com/dayterr/go-diploma/cmd/gophermart/handlers"
)

type Storager interface {
	AddUser(user handlers.UserModel) (int64, error)
}

type UserStorage struct {
	DB *sql.DB
	DSN string
}
