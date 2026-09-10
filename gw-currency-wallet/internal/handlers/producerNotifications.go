package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	model "gw-currency-wallet/internal/storages"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

func KafkaProducer(amount int, id int) {

	fmt.Println("ОТПРАВЛЯЕМ В КАФКУ ")

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	writer := &kafka.Writer{
		Addr:         kafka.TCP(kafkaBroker),
		Topic:        "big-money",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}
	defer writer.Close()
	wg := sync.WaitGroup{}

	var d model.Data

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
	fmt.Println("Сколько вермя прошло", finish)

	// d.Id = id
	// d.Amount = amount

	// data, _ := json.Marshal(d)

	// msg := kafka.Message{
	// 	Value: data,
	// }

	// err := writer.WriteMessages(context.Background(), msg)
	// if err != nil {
	// 	fmt.Println("ОШИБКА ПРИ Оправки в кафку", err)
	// 	return
	// }
	// fmt.Println("Отправлено сообщение ", d.Id, d.Amount)
}
