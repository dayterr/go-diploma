package storage

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func NewUserDB(dsn string) (UserStorage, error) {

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

	err = res.Scan(&order)
	if err != nil {
		log.Println("scanning OrderModel error", err)
		return OrderModel{}, err
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
		log.Println("user creation error:", err)
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, err
}
