package handlers

import (
	"github.com/dgrijalva/jwt-go/v4"
)

type UserIDKey string

type User struct {
	Name string `json:"login"`
	Password string `json:"password"`
}

type Auth struct {
	Key string
}

type CustomClaims struct {
	UserID int64
	jwt.StandardClaims
}