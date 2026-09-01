package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	model "gw-currency-wallet/internal/storages"

	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func BalanceFunc(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		var d model.Data
		err := json.NewDecoder(r.Body).Decode(&d)
		if err != nil {
			fmt.Println("ОШИБКА ПРИ декодировании JSON", err)
			return
		}

		var count int
		err = db.QueryRow(context.Background(),
			"SELECT COUNT(*) FROM operations WHERE user_id = $1",
			d.Id).Scan(&count)
		if err != nil {
			fmt.Println("ОШИБКА ПРИ проверки на наличие Юзера", err)
		}

		if count == 0 {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"error":   "User not found",
				"user_id": d.Id,
			})
			return
		} else {

			var balance int

			err = db.QueryRow(context.Background(),
				"SELECT COALESCE(SUM(CASE WHEN operation = 'deposit' THEN amount WHEN operation = 'withdraw' THEN -amount END), 0) FROM operations WHERE user_id = $1",
				d.Id).Scan(&balance)
			if err != nil {
				fmt.Println("ОШИБКА ПРИ ПОЛУЧЕНИИ ДАННЫХ из БД", err)
			}

			fmt.Println("ПОЛУченныый баланс", balance)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(balance)

		}

	}
}
