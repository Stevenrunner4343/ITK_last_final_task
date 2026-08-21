package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/storages/model"

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
		json.NewDecoder(r.Body).Decode(&a)

		hash := sha256.Sum256([]byte(a.Password))
		password := fmt.Sprintf("%x", hash)

		err := db.QueryRow(context.Background(), "select username FROM accounts where username = $1",
			a.Username).Scan(&username)

		if err == pgx.ErrNoRows { //можжет стоит проверят username а не err

			err1 := db.QueryRow(context.Background(), "select username FROM accounts where email = $1",
				a.Email).Scan(&email)

			if err1 == pgx.ErrNoRows {

				fmt.Println("НОвый Юзер добален")

				db.Exec(context.Background(), "INSERT INTO accounts (username,password,email) VALUES ($1,$2,$3)", a.Username, password, a.Email)

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"email": a.Email, // спрсить про то как без доп запроса к БД получать именно user_id
					"exp":   time.Now().Add(24 * time.Hour).Unix(),
				})
				tokenString, _ := token.SignedString(config.JwtSecret)

				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "User registered successfully",
					"token":   tokenString,
				})

			} else {
				fmt.Println("Имэел занят")

				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Email already used bUT the username is not taken",
				})

			}
		} else {
			fmt.Println("Юзер занят", err, username)

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Username already exists",
			})

		}
	}
}

func SingInFunc(db *pgxpool.Pool) http.HandlerFunc { // напиcать на фронте что бы направляло на SingIn после регистрации
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

				w.WriteHeader(http.StatusBadRequest) //подучить http коды
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
