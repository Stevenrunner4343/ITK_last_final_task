package main

import (
	"context"
	"gw-notification/internal/consumer"
	"gw-notification/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Test struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func main() {

	logger.LoggerInit("gw-notification")
	defer logger.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer cancel()
	consumer.NewConsumer(ctx)

	stop := make(chan os.Signal, 1)

	<-stop
	time.Sleep(2 * time.Second)
	logger.Info(" gracefully sHUT Down complete ")

}
