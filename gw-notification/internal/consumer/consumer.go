package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"gw-notification/internal/storages/model"
	"gw-notification/pkg/logger"
	"os"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

type Consumer struct {
	reader     *kafka.Reader
	client     *mongo.Client
	collection *mongo.Collection
}

func NewConsumer(ctx context.Context) *Consumer {
	c := &Consumer{}

	uri := "mongodb://admin:pass@mongodb:27017/mydb?authSource=admin"

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		logger.Error("Ошибка при создании Клиента Монго", zap.Error(err))
		os.Exit(1)
	}
	c.collection = client.Database("mydb").Collection("testTable")
	c.client = client

	kafkaBroker := os.Getenv("KAFKA_BROKER")

	c.reader = kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:        []string{kafkaBroker},
			Topic:          "big-money",
			GroupID:        "big-money-group",
			StartOffset:    kafka.FirstOffset,
			CommitInterval: 0,
		})

	go c.SendToMongo(ctx)

	return c

}

func (c *Consumer) SendToMongo(ctx context.Context) {

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			logger.Error("ОШИБКА ПРи чтении в очередь Кафка", zap.Error(err))
			continue
		}

		var d model.Data
		err = json.Unmarshal(msg.Value, &d)
		if err != nil {
			logger.Error("ОШИБКА ДЕкодировании JsoN пЕРЕД ВСТАВКОЙ В moNGO", zap.Error(err))
		}
		_, err = c.collection.InsertOne(ctx, d)
		if err != nil {
			logger.Error("ОШИБКА ПРи вставке В МОНго", zap.Error(err))
		}
	}
}
