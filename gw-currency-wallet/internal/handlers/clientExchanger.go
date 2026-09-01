package handlers

import (
	"context"
	"fmt"
	"os"

	pb "github.com/Stevenrunner4343/proto-exchange/exchange"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("ОШИБКА ПРИ ПОДКЛЮЧЕНИИ К GRPC:", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewExchangeServiceClient(conn)

	Rates, err := client.GetExchangeRates(context.Background(), &pb.Empty{})
	if err != nil {
		fmt.Println("ОШИБКА ПРИ получении Курсов!", err)
		os.Exit(1)
	}
	fmt.Println("Все курсы:", Rates.Rates)

	Rate, err := client.GetExchangeRateForCurrency(context.Background(), &pb.CurrencyRequest{
		FromCurrency: "USD",
		ToCurrency:   "RUB",
	})
	if err != nil {
		fmt.Println("ОШИБКА ПРИ получении одного кура", err)
		os.Exit(1)
	}
	fmt.Println("из чего во что и", Rate.FromCurrency, Rate.ToCurrency, Rate.Rate)
}
