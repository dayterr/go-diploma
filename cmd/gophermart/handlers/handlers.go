package handlers

import (
	"github.com/dayterr/go-diploma/internal/storage"
	"io"
	"log"
	"net/http"
	"encoding/json"
	"strconv"
)

type AsyncHandler struct{
	Auth Auth
}

func NewAsyncHandler(dsn string) AsyncHandler {
	var auth Auth
	var s storage.Storager
	s, err := storage.NewUserDB(dsn)
	if err != nil {
		log.Println("setting database error", err)
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
}

func (ah *AsyncHandler) LogUser(w http.ResponseWriter, r *http.Request) {
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

	token, err := ah.Auth.LogUser(u, "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "Bearer",
		Value: token,
	})
	w.WriteHeader(http.StatusOK)
}

func (ah *AsyncHandler) LoadOrderNumber(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orderNumber, err := strconv.Atoi(string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	order, err := ah.Auth.Storage.GetOrder(orderNumber)

	username := r.Context().Value("username").(string)
	userID, err := ah.Auth.Storage.GetUser(username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if int(userID) == order.UserID {
		w.WriteHeader(http.StatusOK)
	} else if order.UserID != 0 {
		w.WriteHeader(http.StatusConflict)
	}

	_, err = ah.Auth.Storage.AddOrder(orderNumber, int(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)

}


