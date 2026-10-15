package consumer

import (
	"context"
	"errors"
	"os"

	"gw-notification/internal/service"
	"gw-notification/pkg/logger"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Consumer struct {
	reader *kafka.Reader
	svc    *service.NotificationService
}

func NewConsumer(svc *service.NotificationService) *Consumer {
	c := &Consumer{svc: svc}

	kafkaBroker := os.Getenv("KAFKA_BROKER")

	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaBroker},
		Topic:          "big-money",
		GroupID:        "big-money-group",
		StartOffset:    kafka.FirstOffset,
		CommitInterval: 0,
	})

	return c
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			logger.Error("Ошиька при Чтении Сообщений из кафки", zap.Error(err))
			continue
		}

		err = c.svc.ProcessMessage(ctx, msg.Value)
		if err != nil {
			logger.Error("Внутення ошибка из Сервиса", zap.Error(err))
			continue
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
