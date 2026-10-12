package handlers

import (
	"encoding/json"
	"errors"
	"gw-currency-wallet/internal/repository"
	"gw-currency-wallet/internal/service"
	model "gw-currency-wallet/internal/storages/model"
	"gw-currency-wallet/pkg/logger"

	"net/http"

	"go.uber.org/zap"
)

func BalanceFunc(s *service.BalanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var d model.Data
		err := json.NewDecoder(r.Body).Decode(&d)
		if err != nil {
			logger.Error("Ошибка при декодировании", zap.Error(err))
			return
		}
		balance, err := s.GetBalance(r.Context(), d.Id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				http.Error(w, "ЮЗер не Найден", http.StatusInternalServerError)
				return
			} else {
				http.Error(w, "Внутрення Ощибка Из Сервиса", http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(balance)

	}

}
