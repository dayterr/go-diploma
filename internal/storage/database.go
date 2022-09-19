package storage

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func NewDB(dsn string) (UserStorage, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return UserStorage{}, err
	}

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS users (id serial PRIMARY KEY, 
                                 login text UNIQUE NOT NULL, 
                                 password text NOT NULL);`)
	if err != nil {
		return UserStorage{}, err
	}
	log.Println("table users created")

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS orders (id serial PRIMARY KEY, 
                                  number BIGINT UNIQUE NOT NULL, 
                                  status TEXT NOT NULL,
                                  accrual DOUBLE PRECISION NOT NULL,
                                  uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL,
                                  user_id INT NOT NULL,
                                  FOREIGN KEY (user_id) REFERENCES public.users(id));`)
	if err != nil {
		return UserStorage{}, err
	}
	log.Println("table orders created")

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS balance (id serial PRIMARY KEY, 
                                   current DOUBLE PRECISION, 
                                   withdrawn DOUBLE PRECISION,
                                   uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL,
                                   user_id INT NOT NULL,
                                   FOREIGN KEY (user_id) REFERENCES public.users(id));`)
	if err != nil {
		return UserStorage{}, err
	}
	log.Println("table balance created")

	_, err = db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS withdrawals (id serial PRIMARY KEY, 
                                       order_number BIGINT UNIQUE NOT NULL,
                                       sum INT NOT NULL, 
                                       uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL,
                                       user_id INT NOT NULL,
                                       FOREIGN KEY (user_id) REFERENCES public.users(id),
    								   FOREIGN KEY (order_number) REFERENCES public.orders(number));`)
	if err != nil {
		return UserStorage{}, err
	}
	log.Println("table withdrawals created")

	return UserStorage{
		DB:           db,
		DSN:          dsn,
	}, nil
}

func (us UserStorage) AddUser(user UserModel) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("writing user to database")
	res, err := us.DB.ExecContext(ctx,
		`INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id`,
		user.Name, user.Password)
	if err != nil {
		log.Println("user creation error:", err)
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, err
}

func (us UserStorage) GetUser(username string) (int64, error) {
	var userID int64

	log.Println("getting user from database")
	res := us.DB.QueryRow(`SELECT id FROM users WHERE login = $1`, username)
	err := res.Scan(&userID)
	if err != nil {
		log.Println("getting id error", err)
		return 0, err
	}
	return userID, nil
}

func (us UserStorage) GetOrder(orderNumber int) (OrderModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var order OrderModel

	log.Println("getting order from database")
	res, err := us.DB.QueryContext(ctx, `SELECT id, number, status, accrual, uploaded_at, user_id 
FROM orders WHERE number = $1`, orderNumber)
	if err != nil {
		log.Println("getting order error", err)
		return OrderModel{}, nil
	}

	for res.Next() {
		err = res.Scan(&order.ID, &order.Number, &order.Status, &order.Accrual,
			&order.UploadedAt, &order.UserID)
		if err != nil {
			log.Println("scanning OrderModel error", err)
			return OrderModel{}, err
		}
	}
	return order, nil
}

func (us UserStorage) AddOrder(orderNumber, userID int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("writing order to database")
	res, err := us.DB.ExecContext(ctx,
		`INSERT INTO orders (number, status, accrual, uploaded_at, user_id) 
                    VALUES ($1, $2, $3, $4, $5) RETURNING id`,
                    orderNumber, "NEW", 0.0, time.Now().Format(time.RFC3339), userID)
	if err != nil {
		log.Println("order creation error:", err)
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, err
}

func (us UserStorage) GetListOrders(userID int) ([]OrderModel, error) {

	log.Println("getting user orders from database")
	res, err := us.DB.Query(`SELECT number, status, accrual, uploaded_at from orders WHERE user_id = $1`, userID)
	if err != nil {
		log.Println("getting orders error", err)
		return []OrderModel{}, err
	}
	defer res.Close()

	var orders []OrderModel
	for res.Next() {
		order := OrderModel{}
		err = res.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			log.Println("scanning order error", err)
			return []OrderModel{}, err
		}
		orders = append(orders, order)
	}

	if res.Err() != nil{
		log.Println("some row error", err)
		return []OrderModel{}, err
	}

	return orders, nil
}

func (us UserStorage) UpdateOrders(order OrderModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("updating orders db")
	_, err := us.DB.ExecContext(ctx, `UPDATE orders SET status = $1, 
                  accrual = $2 WHERE number = $3`, order.Status, order.Accrual, order.Number)
	if err != nil {
		log.Println("updating order error", err)
		return err
	}

	return nil
}

func (us UserStorage) FindUser(orderNumber int) (int64, error) {
	res, err := us.DB.Query(`SELECT user_id FROM orders WHERE number = $1`, orderNumber)
	if err != nil {
		log.Println("finding user error", err)
		return 0, err
	}
	defer res.Close()

	var userID int64
	if res.Next() {
		err = res.Scan(&userID)
		if err != nil {
			log.Println("scanning userID error", err)
			return 0, err
		}
	}

	return userID, nil
}

func (us UserStorage) GetBalance(userID int) (float64, error) {
	res, err := us.DB.Query(`SELECT id, current FROM balance WHERE user_id = $1`, userID)
	if err != nil {
		log.Println("getting balance error", err)
		return 0, err
	}
	defer res.Close()

	var currentBalance float64
	for res.Next() {
		err := res.Scan(&currentBalance)
		if err != nil {
			log.Println("scanning balance error", err)
			return 0, err
		}
	}
	return currentBalance, err
}

func (us UserStorage) UpdateBalance(balance float64, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := us.DB.ExecContext(ctx, `UPDATE balance SET current = $1 WHERE user_id = $2`,
		balance, userID)
	if err != nil {
		log.Println("updating balance error", err)
		return err
	}
	return nil
}

func (us UserStorage) GetFullInfoBalance(userID int) (BalanceModel, error) {
	res, err := us.DB.Query(`SELECT current, withdrawn FROM balance WHERE user_id = $1`, userID)
	if err != nil {
		log.Println("getting balance error", err)
		return BalanceModel{}, err
	}

	var balance BalanceModel
	for res.Next() {
		err = res.Scan(&balance.Current, &balance.Withdrawn)
		if err != nil {
			log.Println("scanning balance error", err)
			return BalanceModel{}, err
		}
	}
	return balance, nil
}