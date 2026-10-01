package main

import (
	"encoding/json"
	"fmt"
	"gw-currency-wallet/internal/handlers"
	"gw-currency-wallet/internal/middleware"
	"gw-currency-wallet/pkg/logger"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

var wg = sync.WaitGroup{}

func produceTest() {

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "kafka:9092"
	}

	writer := &kafka.Writer{
		Addr:  kafka.TCP(kafkaBroker),
		Topic: "analytics",
	}
	defer writer.Close()

	for i := range 10000 {
		wg.Add(1)
		go func() {
			wg.Done()
			data, _ := json.Marshal(map[string]any{
				"id":     i,
				"amount": rand.Intn(70000) + 30000,
			})
			_ = writer.WriteMessages(context.Background(), kafka.Message{Value: data})

		}()

	}
	wg.Wait()
	fmt.Println("10000 сообщений отправлено")
}

func main() {

	logger.LoggerInit("gw-currency-wallet")
	defer logger.Sync()
	produceTest()
	connection := fmt.Sprintf("host=postgres port=5432 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	db, err := pgxpool.New(context.Background(), connection)
	if err != nil {
		fmt.Println("ОШИБКА ПРИ ПОДКЛЮЧЕНИИ К БД", err)
		os.Exit(1)
	}

	producer := handlers.NewProducer()
	defer producer.Close()

	http.HandleFunc("/balance", middleware.Middleware(handlers.BalanceFunc(db)))
	http.HandleFunc("/operation", handlers.OperationFunc(db, producer)) //ОБЕРНУТЬ В MIDDDLEWARE

	http.HandleFunc("/signUp", handlers.SingUpFunc(db))
	http.HandleFunc("/singIn", handlers.SingInFunc(db))

	go func() {
		time.Sleep(10 * time.Second)
		handlers.GetRatesClient()

	}()

	http.ListenAndServe(":8082", nil)

}

/*

SELECT
user_id,
SUM(
CASE
WHEN operation = 'deposit' THEN amount
WHEN operation = 'withdraw' THEN -amount
END
)
FROM operations
GROUP BY user_id
HAVING user_id = 2001




CREATE TABLE operations (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    operation VARCHAR(255),
    amount BIGINT,
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
)
*/
