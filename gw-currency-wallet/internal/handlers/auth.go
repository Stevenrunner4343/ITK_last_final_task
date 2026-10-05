package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gw-currency-wallet/internal/config"
	model "gw-currency-wallet/internal/storages/model"

	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SingUpFunc(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		var a model.Account
		var username string
		var email string
		var userID int
		json.NewDecoder(r.Body).Decode(&a)

		hash := sha256.Sum256([]byte(a.Password))
		password := fmt.Sprintf("%x", hash)

		err := db.QueryRow(context.Background(), "select username FROM accounts where username = $1 AND email = $2",
			a.Username, a.Email).Scan(&username, &email)

		if err == pgx.ErrNoRows {
			// вместо Email сделал user ID и написал один запрос всместо двух! протестировать
			fmt.Println("НОвый Юзер добален")

			db.QueryRow(context.Background(), "INSERT INTO accounts (username,password,email) VALUES ($1,$2,$3) RETURNING id", a.Username, password, a.Email).Scan(&userID)

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userId": userID,
				"exp":    time.Now().Add(24 * time.Hour).Unix(),
			})
			tokenString, _ := token.SignedString(config.JwtSecret)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "User REgistered",
				"token":   tokenString,
			})

		} else {
			fmt.Println("Email или Username зАнят!", err, username, email)

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Username already exists",
			})

		}
	}
}

func SingInFunc(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		var a model.Account
		var username string
		json.NewDecoder(r.Body).Decode(&a)

		hash := sha256.Sum256([]byte(a.Password))
		password := fmt.Sprintf("%x", hash)

		err := db.QueryRow(context.Background(), "SELECT username FROM accounts WHERE username = $1 ",
			a.Username).Scan(&username)
		if err == pgx.ErrNoRows {

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid username",
			})
			fmt.Println("SIGn In нет такого юзера")

		} else {

			err1 := db.QueryRow(context.Background(), "SELECT username FROM accounts WHERE password = $1 ",
				password).Scan(&password)
			if err1 == pgx.ErrNoRows {

				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Invalid password",
				})
				fmt.Println("SIGn In не верный пароль ")

			} else {
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "succsessfully logged",
				})

			}

		}
	}
}
