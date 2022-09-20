package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/dayterr/go-diploma/internal/storage"
	"github.com/dgrijalva/jwt-go/v4"
	"log"
	"net/http"
	"time"
)

func EncryptPassword(password string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func createToken(id int64, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomClaims{
		UserID: id,
		StandardClaims: jwt.StandardClaims {
			ExpiresAt: jwt.At(time.Now().Add(1440 * time.Minute)),
			IssuedAt:  jwt.At(time.Now()),
		},
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
	if err != nil {
		log.Println("adding user error", err)
		return "", err
	}

	token, err := createToken(id, key)
	return token, err
}

func (a Auth) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urls := []string{"/api/user/register", "/api/user/login"}
		path := r.URL.Path
		for _, v := range urls {
			if v == path {
				next.ServeHTTP(w, r)
				return
			}
		}

		cookieToken, err := r.Cookie("Bearer")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims := &CustomClaims{}
		token, err := jwt.ParseWithClaims(cookieToken.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.Key), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey("userid"), claims.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	})
}

func (a Auth) LogUser(user User, key string) (string, error) {
	var modelUser storage.UserModel
	modelUser.Name = user.Name
	//modelUser.Password = EncryptPassword(user.Password, key)

	id, err := a.Storage.GetUser(modelUser.Name)
	if err != nil {
		log.Println("getting user error", err)
		return "", err
	}

	token, err := createToken(id, key)
	return token, err
}
