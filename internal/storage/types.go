package storage

import (
	"database/sql"
)

type Storager interface {
	AddUser(user UserModel) (int64, error)
	GetUser(username string) (int64, error)
	GetOrder(orderNumber int) (OrderModel, error)
	AddOrder(orderNumber, userID int) (int64, error)
	GetListOrders(userID int) ([]OrderModel, error)
	UpdateOrders(order OrderModel) error
	FindUser(orderNumber int) (int64, error)
	GetBalance(userID int) (float64, error)
	UpdateBalance(balance float64, userID int) error
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
	Number int
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