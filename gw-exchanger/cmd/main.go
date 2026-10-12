package main

import (
	"context"
	"gw-exchanger/internal/handlers"
	"gw-exchanger/internal/storages/postgres"
	"gw-exchanger/pkg/logger"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	blocker := make(chan struct{})
	logger.LoggerInit("gw-exchanger")
	defer logger.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	db, err := postgres.NewConnection(ctx)
	if err != nil {
		logger.Error("не удалось подключиться к Postgres", zap.Error(err))
		return
	}
	defer db.Close()

	go handlers.StartGrpcServer(db)
	<-blocker
	// можно ли не блокером это сделать а gracefull shutDown
	// stop := make(chan os.Signal, 1)

	// signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	// <-stop

}
