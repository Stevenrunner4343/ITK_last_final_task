package handlers

import (
	"encoding/json"
	"errors"
	"gw-currency-wallet/internal/repository"
	"gw-currency-wallet/internal/service"
	model "gw-currency-wallet/internal/storages/model"

	"net/http"
)

func SingUpFunc(auth *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var a model.Account

		json.NewDecoder(r.Body).Decode(&a)

		tokenString, err := auth.MakeNewUser(r.Context(), a)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Такой Юзеру уже есть",
				})
				return
			} else {
				http.Error(w, "Внутрення Ошибка Из Сервиса", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Юзер зарегетсрирован",
			"token":   tokenString,
		})

	}
}

func SingInFunc(auth *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var a model.Account
		json.NewDecoder(r.Body).Decode(&a)

		err := auth.CheckUser(r.Context(), a)

		if err != nil {
			if err == repository.ErrNotFound {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "НЕ верное Имя",
				})
				return

			}
			if err == repository.ErrWrongPassword {

				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Не верный пароль",
				})
				return

			} else {
				http.Error(w, "Внутрення Ошибка Из Сервиса", http.StatusInternalServerError)
				return
			}

		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "successfully logged in",
		})
	}

}
