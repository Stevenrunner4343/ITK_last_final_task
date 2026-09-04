package main

import (
	"fmt"
	"gw-currency-wallet/internal/handlers"
	"gw-currency-wallet/internal/middleware"
	"net/http"
	"os"
	"time"

	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	connection := fmt.Sprintf("host=postgres port=5432 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	db, err := pgxpool.New(context.Background(), connection)
	if err != nil {
		fmt.Println("ОШИБКА ПРИ ПОДКЛЮЧЕНИИ К БД", err)
		os.Exit(1)
	}

	http.HandleFunc("/balance", middleware.Middleware(handlers.BalanceFunc(db)))
	http.HandleFunc("/operation", middleware.Middleware(handlers.OperationFunc(db)))

	http.HandleFunc("/signUp", handlers.SingUpFunc(db))
	http.HandleFunc("/singIn", handlers.SingInFunc(db))
	go func() {
		time.Sleep(15 * time.Second)
		handlers.GetRatesClient()

	}()

	fmt.Println("Сервер запущен на 8082")
	http.ListenAndServe(":8082", nil)

}

/*

SELECT
user_id,
SUM(
CASE
WHEN operation = 'deposit' THEN amount
WHEN operation = 'withdraw' THEN -amount
END
)
FROM operations
GROUP BY user_id
HAVING user_id = 2001




CREATE TABLE operations (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    operation VARCHAR(255),
    amount BIGINT,
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
)
*/
