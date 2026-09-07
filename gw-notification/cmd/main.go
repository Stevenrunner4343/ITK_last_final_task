package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var topic string = "big-money"

func Producer() {

	fmt.Println("ОТПРАВЛЯЕМ В КАФКУ ")

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	writer := &kafka.Writer{
		Addr:         kafka.TCP(kafkaBroker),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}
	defer writer.Close()

	for i := range 5 {
		msg := kafka.Message{
			Key:   fmt.Appendf(nil, "key-%d", i),
			Value: fmt.Appendf(nil, `{"amount":%d}`, 3000+i),
		}

		err := writer.WriteMessages(context.Background(), msg)
		if err != nil {
			fmt.Println("ОШИБКА ПРИ Оправки в кафку", err)
			os.Exit(1)
		}
		fmt.Println("Отправлено сообщение ", i)
	}

}

func Consumer() {

	fmt.Println(" Читаем ИЗ кафКИ")
	kafkaBroker := os.Getenv("KAFKA_BROKER")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaBroker},
		Topic:          topic,
		GroupID:        "big-money-group",
		StartOffset:    kafka.FirstOffset,
		CommitInterval: 0,
	})
	defer reader.Close()

	time.Sleep(1 * time.Second)
	for range 10 {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("ОШИБКА ПРИ ЧТЕНИИИ", err)
			os.Exit(1)
		}
		fmt.Println(" Получено:", string(msg.Value))
	}

	fmt.Println(" ВЕС РАБОТАЕТ!!!!")

}

type Test struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func sendToMongo() {
	uri := "mongodb://admin:pass@mongodb:27017/mydb?authSource=admin"

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("ОШИБКА ПРИ пОДКЛЮЧЕНИИ К мОНГО", err)
		os.Exit(1)
	}
	defer client.Disconnect(context.Background())
	collection := client.Database("mydb").Collection("testTable")

	var d Test
	d.Name = "bye"
	d.Age = 26

	result, err := collection.InsertOne(context.Background(), d)

	if err != nil {
		fmt.Println("ОШИБКА ПРИ всатвке", err)
	}
	fmt.Println("резутать при вставке!", result)
}
func main() {
	// Producer()
	// Consumer()
	sendToMongo()

}
