package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/dayterr/go-diploma/internal/storage"
	"github.com/dgrijalva/jwt-go/v4"
	"log"
	"time"
)

func EncryptPassword(password string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func createToken(id int64, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: jwt.At(time.Now().Add(1440 * time.Minute)),
		IssuedAt:  jwt.At(time.Now()),
	})
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		log.Print("signing token error:", err)
		return "", err
	}
	return signedToken, nil
}

func (a Auth) RegisterNewUser(user User, key string) (string, error) {
	var modelUser storage.UserModel
	modelUser.Name = user.Name
	modelUser.Password = EncryptPassword(user.Password, key)

	id, err := a.Storage.AddUser(modelUser)
	token, err := createToken(id, key)
	return token, err
}

func (a Auth) LogUser(user User, key string) (string, error) {
	var modelUser storage.UserModel
	modelUser.Name = user.Name
	modelUser.Password = EncryptPassword(user.Password, key)

	id, err := a.Storage.GetUser(modelUser)
	token, err := createToken(id, key)
	return token, err
}
