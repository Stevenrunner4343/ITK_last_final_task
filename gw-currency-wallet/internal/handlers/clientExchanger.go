package handlers

import (
	"context"
	"fmt"
	"os"

	pb "github.com/Stevenrunner4343/proto-exchange/exchange"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetRatesClient() {
	grpcHost := os.Getenv("GRPC_HOST")
	conn, err := grpc.NewClient(grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("ОШИБКА ПРИ СОЗДАНИИ КЛИЕНТ Grpc:", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewExchangeServiceClient(conn)

	Rates, err := client.GetExchangeRates(context.Background(), &pb.Empty{})
	if err != nil {
		fmt.Println("ОШИБКА ПРИ ПОДКЛЮЧЕНИИ К GRPC", err)
		os.Exit(1)
	}
	fmt.Println("Все курсы:", Rates.Rates)

	Rate, err := client.GetExchangeRateForCurrency(context.Background(), &pb.CurrencyRequest{
		FromCurrency: "USD",
		ToCurrency:   "RUB",
	})
	if err != nil {
		fmt.Println("ОШИБКА  ПРИ получении одного кура", err)
		os.Exit(1)
	}
	fmt.Println("из чего во что ", Rate.FromCurrency, Rate.ToCurrency, Rate.Rate)

}

func main() {

}
