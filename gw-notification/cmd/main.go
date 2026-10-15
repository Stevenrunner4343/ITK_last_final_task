package main

import (
	"context"
	"os/signal"
	"syscall"

	"gw-notification/internal/consumer"
	"gw-notification/internal/repository"
	"gw-notification/internal/service"
	"gw-notification/pkg/logger"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

func main() {
	logger.LoggerInit("gw-notification")
	defer logger.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	uri := "mongodb://admin:pass@mongodb:27017/mydb?authSource=admin"
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		logger.Error("ОШибка при подключении к монго!", zap.Error(err))
		return
	}
	defer client.Disconnect(context.Background())

	repo := repository.NewNotificationRepo(client, "mydb", "testTable")

	svc := service.NewNotificationService(repo)

	cons := consumer.NewConsumer(svc)
	defer cons.Close()

	go cons.Run(ctx)

	<-ctx.Done()

}
