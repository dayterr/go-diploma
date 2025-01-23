package handlers

import (
	"encoding/json"
	"github.com/dayterr/go-diploma/internal/accrual"
	"github.com/dayterr/go-diploma/internal/storage"
	"io"
	"log"
	"net/http"
	"strconv"
)

type AsyncHandler struct{
	Auth Auth
	AccrualClient accrual.AccrualClient
	Storage storage.Storager
}

type UserID string

func NewAsyncHandler() *AsyncHandler {
	var auth Auth
	auth.Key = ""
	ah := AsyncHandler{Auth: auth}
	return &ah
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

	token, err := ah.RegisterNewUser(u, "")
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

	token, err := ah.LogUserWithToken(u, "")
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

func (ah *AsyncHandler) PostOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	num, err := strconv.Atoi(string(body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !CheckLuhn(num) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	orderNumber := string(body)

	order, err := ah.Storage.GetOrder(orderNumber)
	if err != nil {
		log.Println("getting order error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value(UserIDKey("userid")).(int64)

	if int(userID) == order.UserID {
		w.WriteHeader(http.StatusOK)
		return
	} else if order.UserID != 0 {
		w.WriteHeader(http.StatusConflict)
		return
	}

	_, err = ah.Storage.AddOrder(orderNumber, int(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("order saved")
	log.Println(ah.AccrualClient.OrderChannel)
	ah.AccrualClient.ManagePoints(orderNumber)
	//ah.AccrualClient.OrderChannel <- orderNumber
	log.Println("order added")

	w.WriteHeader(http.StatusAccepted)
}

func (ah *AsyncHandler) LoadOrderList(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey("userid")).(int64)

	orders, err := ah.Storage.GetListOrders(int(userID))
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

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	w.WriteHeader(http.StatusOK)
}

func (ah *AsyncHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey("userid")).(int64)

	balance, err := ah.Storage.GetFullInfoBalance(int(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	body, err := json.Marshal(&balance)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	w.WriteHeader(http.StatusOK)
}

func (ah *AsyncHandler) WithdrawPoints(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey("userid")).(int64)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var withdraw storage.Withdraw

	err = json.Unmarshal(body, &withdraw)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	balance, err := ah.Storage.GetBalance(int(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if withdraw.Sum > balance {
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

	diff := balance - withdraw.Sum
	err = ah.Storage.UpdateBalance(diff, withdraw.Sum, int(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = ah.Storage.AddWithdrawal(withdraw.Sum, withdraw.Order, int(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ah *AsyncHandler) GetListWithdrawal(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey("userid")).(int64)

	withdrawals, err := ah.Storage.GetListWithdrawal(int(userID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	body, err := json.Marshal(&withdrawals)
	if err != err {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	w.WriteHeader(http.StatusOK)

}