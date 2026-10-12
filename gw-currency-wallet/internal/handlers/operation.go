package handlers

import (
	"encoding/json"

	"gw-currency-wallet/internal/service"
	model "gw-currency-wallet/internal/storages/model"
	"gw-currency-wallet/pkg/logger"

	"net/http"

	"go.uber.org/zap"
)

func OperationFunc(s *service.OperationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var d model.Data

		err := json.NewDecoder(r.Body).Decode(&d)
		if err != nil {
			logger.Error("Ошибка при Декодирвоание Данных Из ЗАпроса", zap.Error(err), zap.Int("user_id", d.Id))
			return
		}

		err = s.Process(r.Context(), d)
		if err != nil {
			http.Error(w, "Внутрення Ощибка Из Сервиса", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
