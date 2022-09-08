package handlers

import (
	"github.com/dayterr/go-diploma/internal/storage"
	"io"
	"log"
	"net/http"
	"encoding/json"
)

type AsyncHandler struct{
	Auth Auth
}

func NewAsyncHandler(dsn string) AsyncHandler {
	var auth Auth
	var s storage.Storager
	s, err := storage.NewUserDB(dsn)
	if err != nil {
		log.Println(err)
	}
	auth.Storage = s
	auth.Key = ""
	ah := AsyncHandler{Auth: auth}
	return ah
}

func (ah *AsyncHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var u User
	err = json.Unmarshal(body, &u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	token, err := ah.Auth.RegisterNewUser(u, "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Bearer",
		Value: token,
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("JWT " + token))
}


