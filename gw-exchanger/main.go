package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ExchangeFunc(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

	}
}

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

	http.HandleFunc("/exchange", ExchangeFunc(db))

}
