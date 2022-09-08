package storage

import "database/sql"

type UserStorage struct {
	DB *sql.DB
	DSN string
	Login string
	Password string
}
