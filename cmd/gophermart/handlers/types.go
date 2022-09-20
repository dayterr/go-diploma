package handlers

import (
	"github.com/dayterr/go-diploma/internal/storage"
	"github.com/dgrijalva/jwt-go/v4"
)

type UserIDKey string

type User struct {
	Name string `json:"login"`
	Password string `json:"password"`
}

type Auth struct {
	Key string
	Storage storage.Storager
}

type CustomClaims struct {
	UserID int64
	jwt.StandardClaims
}