package consumer

import (
	"context"
	"errors"
	"gw-analytics/pkg/logger"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Consumer struct {
	reader *kafka.Reader
	conn   driver.Conn
}

func NewConsumer(ctx context.Context, chConn driver.Conn) *Consumer {
	c := &Consumer{}
	c.conn = chConn

	kafkaBroker := os.Getenv("KAFKA_BROKER")

	c.reader = kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:        []string{kafkaBroker},
			Topic:          "analytics",
			GroupID:        "analytics-group",
			StartOffset:    kafka.FirstOffset,
			CommitInterval: 0,
		})

	// go func() {
	// 	for {
	// 		msg, err := c.reader.ReadMessage(ctx)

	// 		if err != nil {
	// 			if errors.Is(err, context.Canceled) {
	// 				logger.Info("Отменили КОнтекст Отснавливаемся!")
	// 				return
	// 			}
	// 			logger.Error("Ошибка чтенияСообщения из КОнсюмера", zap.Error(err))
	// 		}

	// 		err = c.conn.Exec(ctx,
	// 			"INSERT INTO events (ts, event_type, status) VALUES (?, ?, ?)",
	// 			time.Now(), "message", string(msg.Value),
	// 		)
	// 		if err != nil {
	// 			logger.Error("Ошибка вставки в ClickHouse", zap.Error(err))
	// 			return
	// 		}
	// 		fmt.Println("вставлено ", string(msg.Value))

	// 	}

	// }()
	go c.run(ctx)

	return c

}
func (c *Consumer) run(ctx context.Context) {
	var buffer []kafka.Message
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	flush := func() {
		if len(buffer) == 0 {
			return
		}
		batch, err := c.conn.PrepareBatch(ctx, "INSERT INTO events (ts, event_type, status)")
		if err != nil {
			logger.Error("PrepareBatch", zap.Error(err))
			return
		}
		for _, m := range buffer {
			_ = batch.Append(time.Now(), "message", string(m.Value))
		}
		if err := batch.Send(); err != nil {
			logger.Error("batch.Send", zap.Error(err))
			return
		}
		logger.Info("batch sent", zap.Int("rows", len(buffer)))
		buffer = buffer[:0] // переиспользуем слайс
	}

	for {
		select {
		case <-ticker.C:
			flush()
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				logger.Error("read", zap.Error(err))
				continue
			}
			buffer = append(buffer, msg)
			if len(buffer) >= 5000 {
				flush()
			}
		}
	}
}
func (c *Consumer) Close() {
	err := c.reader.Close()
	if err != nil {
		logger.Error("Ошибка При закрытии Консюмера! ", zap.Error(err))
	}
	// может и не надо логировать закртие вообще
}
