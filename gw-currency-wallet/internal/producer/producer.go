package producer

import (
	"context"
	"encoding/json"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/pkg/logger"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const (
	TopicAnalytics = "analytics"
	TopicBigMoney  = "big-money"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer() *Producer {
	p := &Producer{}
	kafkaBroker := config.KafkaBroker
	p.writer = &kafka.Writer{
		Addr:         kafka.TCP(kafkaBroker),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}
	return p
}

func (p *Producer) Send(ctx context.Context, data any, topic string) error {

	result, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Topic: topic,
		Value: result,
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		return err
	}
	return nil

}
func (p *Producer) Close() {
	err := p.writer.Close()
	if err != nil {
		logger.Error("Ошибка при Закрытии Продсера", zap.Error(err))
	}
}
