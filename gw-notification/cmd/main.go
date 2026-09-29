package main

import (
	"context"
	"fmt"
	"gw-notification/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Consumer() {

	fmt.Println(" Читаем ИЗ кафКИ")
	kafkaBroker := os.Getenv("KAFKA_BROKER")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaBroker},
		Topic:          "big-money",
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

	logger.LoggerInit("gw-notification")
	defer logger.Sync()

	logger.Info("Запуск: gw-notification")

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	time.Sleep(2 * time.Second)
	logger.Info(" gracefully sHUT Down complete ")

	// time.Sleep(10 * time.Second)
	Consumer()
	sendToMongo()

}
