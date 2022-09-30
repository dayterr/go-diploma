package storage

import (
	"database/sql"
)

type Storager interface {
	AddUser(user User) (int64, error)
	GetUser(username string) (int64, error)
	GetOrder(orderNumber string) (Order, error)
	AddOrder(orderNumber string, userID int) (int64, error)
	GetListOrders(userID int) ([]Order, error)
	UpdateOrder(order Order) error
	FindUser(orderNumber string) (int64, error)
	GetBalance(userID int) (float64, error)
	UpdateBalance(balance, withdrawn float64, userID int) error
	GetFullInfoBalance(userID int) (Balance, error)
	AddWithdrawal(withdrawn float64, orderNumber string, userID int) error
	GetListWithdrawal(userID int) ([]Withdraw, error)
}

type UserStorage struct {
	DB *sql.DB
	DSN string
}

type User struct {
	ID int
	Name string
	Password string
}

type Order struct {
	ID int
	Number string
	Status string
	Accrual float64
	UploadedAt string
	UserID int
}

type Balance struct {
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