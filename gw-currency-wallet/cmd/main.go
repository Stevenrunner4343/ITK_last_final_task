package main

import (
	"context"
	"gw-currency-wallet/internal/handlers"
	"gw-currency-wallet/internal/middleware"
	"gw-currency-wallet/internal/producer"
	"gw-currency-wallet/internal/repository"
	"gw-currency-wallet/internal/service"
	"gw-currency-wallet/internal/storages/postgres"
	"gw-currency-wallet/pkg/logger"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {

	logger.LoggerInit("gw-currency-wallet")
	defer logger.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	db, err := postgres.NewConnection(ctx)
	if err != nil {
		logger.Error("не удалось подключиться к Postgres", zap.Error(err))
		return
	}
	repo := repository.NewWalletRepo(db)
	defer db.Close()
	prod := producer.NewProducer()
	defer prod.Close()

	OperationServ := service.NewOperationService(repo, prod)
	BalanceServ := service.NewBalanceService(repo)
	AuthServ := service.NewAuthService(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /balance", middleware.Middleware(handlers.BalanceFunc(BalanceServ)))
	mux.HandleFunc("POST /operation", middleware.Middleware(handlers.OperationFunc(OperationServ)))
	mux.HandleFunc("POST /signUp", handlers.SingUpFunc(AuthServ))
	mux.HandleFunc("POST /singIn", handlers.SingInFunc(AuthServ))

	if err := http.ListenAndServe(":8082", mux); err != nil {
		logger.Error("Ошибка При запуске Http сервера", zap.Error(err))
	}

	go func() {
		time.Sleep(10 * time.Second)
		handlers.GetRatesClient()

	}()

}
