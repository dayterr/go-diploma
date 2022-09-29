package storage

import (
	"database/sql"
)

type Storager interface {
	AddUser(user UserModel) (int64, error)
	GetUser(username string) (int64, error)
	GetOrder(orderNumber string) (OrderModel, error)
	AddOrder(orderNumber string, userID int) (int64, error)
	GetListOrders(userID int) ([]OrderModel, error)
	UpdateOrders(order OrderModel) error
	FindUser(orderNumber string) (int64, error)
	GetBalance(userID int) (float64, error)
	UpdateBalance(balance, withdrawn float64, userID int) error
	GetFullInfoBalance(userID int) (BalanceModel, error)
	AddWithdrawal(withdrawn float64, orderNumber string, userID int) error
	GetListWithdrawal(userID int) ([]Withdraw, error)
}

type UserStorage struct {
	DB *sql.DB
	DSN string
}

type UserModel struct {
	ID int
	Name string
	Password string
}

type OrderModel struct {
	ID int
	Number string
	Status string
	Accrual float64
	UploadedAt string
	UserID int
}

type BalanceModel struct {
	ID int
	Current float64
	Withdrawn float64
}

type Withdraw struct {
	ID int
	Order string
	Sum float64
	ProcessedAt string
	UserID int
}