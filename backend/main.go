package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"context"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Data struct {
	Id        int
	Operation string
	Amount    int
}

type Account struct {
	Username string
	Password string
	Email    string
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ") //bearer тип токена
		if token == "" {
			fmt.Println("тОКЕНТ пустой")
			return
		}

		tokenCheck, _ := jwt.Parse(token, func(t *jwt.Token) (any, error) {
			return jwtSecret, nil
		})
		if !tokenCheck.Valid {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			fmt.Println("ОШИБКА ПРИ Расшифровки Токена ")
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", tokenCheck.Claims.(jwt.MapClaims)["user_id"])

		next(w, r.WithContext(ctx))
	}

}

func operationFunc(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		var d Data
		json.NewDecoder(r.Body).Decode(&d)

		_, err := db.Exec(context.Background(),
			"INSERT INTO operations (user_id, operation, amount) VALUES ($1, $2, $3)",
			d.Id, d.Operation, d.Amount)
		if err != nil {
			fmt.Println("ОШИБКА ПРИ Зипис в бд", err)
		}

	}
}

func singUpFunc(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		var a Account
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
					"user_id": a.Username,
					"exp":     time.Now().Add(24 * time.Hour).Unix(),
				})
				tokenString, _ := token.SignedString(jwtSecret)

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

func singInFunc(db *pgxpool.Pool) http.HandlerFunc { // напиать нп фронте что бы  направляло на SingIn после регитсрации
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		var a Account
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

				w.WriteHeader(http.StatusBadRequest) // надо ли резобратся как именно и что тут писать и какие коды и  как их именн отправлять?

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

func balanceFunc(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		var d Data
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

	http.HandleFunc("/balance", Middleware(balanceFunc(db)))
	http.HandleFunc("/operation", Middleware(operationFunc(db)))

	http.HandleFunc("/signUp", singUpFunc(db))
	http.HandleFunc("/singIn", singInFunc(db))

	fmt.Println("Сервер запущен на 8082")
	http.ListenAndServe(":8082", nil)

}

/*

SELECT
user_id,
SUM(
CASE
WHEN operation = 'deposit' THEN amount
WHEN operation = 'withdraw' THEN -amount
END
)
FROM operations
GROUP BY user_id
HAVING user_id = 2001




CREATE TABLE operations (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    operation VARCHAR(255),
    amount BIGINT,
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
)
*/
