package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	model "gw-currency-wallet/internal/storages"

	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func OperationFunc(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		var d model.Data
		json.NewDecoder(r.Body).Decode(&d)
		if d.Amount >= 30000 {
			KafkaProducer(d.Amount, d.Id)
			fmt.Println("ПОлучили большк 30к")
		}
		_, err := db.Exec(context.Background(),
			"INSERT INTO operations (user_id, operation, amount) VALUES ($1, $2, $3)",
			d.Id, d.Operation, d.Amount)
		if err != nil {
			fmt.Println("ОШИБКА ПРИ Зипис в бд", err)
		}

	}
}
