package middleware

import (
	"context"
	"gw-currency-wallet/internal/config"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			http.Error(w, `{"error":"no token"}`, http.StatusUnauthorized)
			return
		}

		tokenCheck, _ := jwt.Parse(token, func(t *jwt.Token) (any, error) {
			return config.JwtSecret, nil
		})

		if !tokenCheck.Valid {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		claims := tokenCheck.Claims.(jwt.MapClaims)
		email := claims["email"].(string)

		ctx := context.WithValue(r.Context(), "email", email)
		next(w, r.WithContext(ctx))
	}

}
