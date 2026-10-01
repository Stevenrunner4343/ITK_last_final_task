package handlers

import (
	"context"
	"encoding/json"
	model "gw-currency-wallet/internal/storages"
	"gw-currency-wallet/pkg/logger"
	"os"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
}

const (
	TopicBigMoney  = "big-money"
	TopicAnalytics = "analytics"
)

func NewProducer() *Producer {
	p := &Producer{}
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	p.writer = &kafka.Writer{
		Addr: kafka.TCP(kafkaBroker),
		// Topic:        topic, не нужен потомучто Потому что передаем в Сообщении
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}
	return p
}

func (p *Producer) Send(ctx context.Context, data model.Data, topic string) error {
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
		logger.Error("Ошибка при Закрытии Продсера", zap.Error(err)) //может это и не надо логировать
		return
	}
}

/*

Тест нагрузки!

wg := sync.WaitGroup{}

	start := time.Now()
	for i := range 1000 {
		wg.Add(1)

		go func(num int) {
			defer wg.Done()
			n := rand.Intn(70000) + 30000

			d.Id = num
			d.Amount = n

			data, _ := json.Marshal(d)

			msg := kafka.Message{
				Value: data,
			}

			fmt.Println(num)

			err := writer.WriteMessages(context.Background(), msg)
			if err != nil {
				fmt.Println("ОШИБКА ПРИ Оправки в кафку", err)
				return
			}

		}(i)
	}
	wg.Wait()

	finish := time.Since(start)
	fmt.Println("Сколько вермя прошло", finish)*/
