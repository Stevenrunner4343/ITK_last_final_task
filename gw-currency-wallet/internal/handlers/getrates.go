package handlers

import (
	"encoding/json"
	"gw-currency-wallet/internal/service"
	"gw-currency-wallet/internal/storages/model"
	"net/http"
)

func GetRates(s *service.RatesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		result, err := s.Rates(r.Context())
		if err != nil {
			http.Error(w, "Внутрення Ошибка Из Сервиса", http.StatusInternalServerError)
			return

		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)

	}

}
func GetRate(s *service.RatesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var c model.Currency
		err := json.NewDecoder(r.Body).Decode(&c)

		result, err := s.Rate(r.Context(), c.From, c.To)
		if err != nil {
			http.Error(w, "Внутрення Ошибка Из Сервиса", http.StatusInternalServerError)
			return

		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)

	}

}
