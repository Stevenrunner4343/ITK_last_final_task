package handlers

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	model "gw-currency-wallet/internal/storages/model"
	"gw-currency-wallet/pkg/logger"

	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func OperationFunc(db *pgxpool.Pool, producer *Producer) http.HandlerFunc {
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
			logger.Error("Ошибка при Декодирвоание Данных Из ЗАпроса", zap.Error(err), zap.Int("user_id", d.Id))
			return
		}
		if d.Amount >= 30000 {
			err := producer.Send(r.Context(), d, TopicBigMoney)
			if err != nil {
				logger.Error("Ошибка при отправки Брокером сообщения о Боольшой сумме в Notification", zap.Error(err), zap.Int("user_id", d.Id))
				return
			}
		}

		_, err = db.Exec(context.Background(),
			"INSERT INTO operations (user_id, operation, amount) VALUES ($1, $2, $3)",
			d.Id, d.Operation, d.Amount)
		if err != nil {
			logger.Error("ОШИБКА ЗАПИСИ в бд", zap.Error(err), zap.Int("user_id", d.Id))
			return
		}
		var e model.Event
		e.ID = rand.Uint64() // число с 19 нулями коллизия будет через 1000 лет
		e.Operation = d.Operation
		e.UserID = d.Id
		e.Amount = d.Amount
		e.CreatedAt = time.Now()

		err = producer.Send(r.Context(), e, TopicAnalytics)
		if err != nil {
			logger.Error("Ошибка при отправки Брокером сообщения в Аналитику", zap.Error(err), zap.Int("user_id", d.Id))
			return
		}

	}
}
