package handlers

import "github.com/dayterr/go-diploma/internal/storage"

type User struct {
	Name string `json:"login"`
	Password string `json:"password"`
}

type Auth struct {
	Key string
	Storage storage.Storager
}