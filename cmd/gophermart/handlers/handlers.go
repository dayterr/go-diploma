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

type UserID string

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
		return
	}

	/*if !CheckLuhn(orderNumber) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}*/

	order, err := ah.Auth.Storage.GetOrder(orderNumber)

	log.Println(order)

	userID := r.Context().Value("userid").(int64)

	if int(userID) == order.UserID {
		w.WriteHeader(http.StatusOK)
		return
	} else if order.UserID != 0 {
		w.WriteHeader(http.StatusConflict)
		return
	}

	_, err = ah.Auth.Storage.AddOrder(orderNumber, int(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (ah AsyncHandler) LoadOrderList(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userid").(int64)

	orders, err := ah.Auth.Storage.GetListOrders(int(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	body, err := json.Marshal(&orders)
	if err != err {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
	w.Header().Set("Content-Type", "application/json")

}


