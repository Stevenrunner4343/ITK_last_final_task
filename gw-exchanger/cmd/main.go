package main

import (
	"context"
	"net"
	"os/signal"
	"syscall"

	"gw-exchanger/internal/handlers"
	"gw-exchanger/internal/repository"
	"gw-exchanger/internal/service"
	"gw-exchanger/internal/storages/postgres"
	"gw-exchanger/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	logger.LoggerInit("gw-exchanger")
	defer logger.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	db, err := postgres.NewConnection(ctx)
	if err != nil {
		logger.Error("postgres", zap.Error(err))
		return
	}
	defer db.Close()

	repo := repository.NewExchangeRepo(db)

	svc := service.NewExchangeService(repo)

	grpcServer := handlers.StartGrpcServer(svc)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("listen", zap.Error(err))
		return
	}

	go func() {
		logger.Info("gRPC Сервер Запустился /main")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("Ошибка во время работы Сервера /main", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down...")

	grpcServer.GracefulStop()
}
