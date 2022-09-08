package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func EncryptPassword(password string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func (a Auth) RegisterNewUser(user User, key string) (string, error) {
	var modelUser UserModel
	modelUser.Name = user.Name
	modelUser.Password = EncryptPassword(user.Password, key)


}
