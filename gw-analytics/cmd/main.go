package main

import (
	"context"

	"gw-analytics/internal/consumer"
	"gw-analytics/internal/storage/clickhouse"
	"gw-analytics/pkg/logger"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	logger.LoggerInit("gw-analytics")
	defer logger.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer cancel()

	chConn, err := clickhouse.NewConnection()
	if err != nil {
		logger.Error("Ошибка при подключении к ClickHouse", zap.Error(err))
	}
	defer chConn.Close()
	logger.Info("Соединение с ClickHouse установлено")

	reader := consumer.NewConsumer(ctx, chConn)
	defer reader.Close()
	<-ctx.Done()

}
