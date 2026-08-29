package main

// gw-exchanger/cmd/main.go

import (
	"context"
	"fmt"
	"gw-exchanger/internal/storages"
	pb "gw-exchanger/proto/exchange"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type ExchangeServer struct {
	pb.UnimplementedExchangeServiceServer
	db *pgxpool.Pool
}

func (s *ExchangeServer) GetExchangeRates(ctx context.Context, empty *pb.Empty) (*pb.ExchangeRatesResponse, error) {
	var c storages.Currency

	err := s.db.QueryRow(ctx, "SELECT * FROM currency").
		Scan(&c.USD, &c.EUR, &c.RUB)
	if err != nil {
		return nil, err
	}

	return &pb.ExchangeRatesResponse{
		Rates: map[string]float32{
			"USD": c.USD,
			"EUR": c.EUR,
			"RUB": c.RUB,
		},
	}, nil
}

func (s *ExchangeServer) GetExchangeRateForCurrency(ctx context.Context, req *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	var c storages.Currency

	err := s.db.QueryRow(ctx, "SELECT * FROM currency").
		Scan(&c.USD, &c.EUR, &c.RUB)
	if err != nil {
		return nil, err
	}

	var rate float32

	if req.ToCurrency == "USD" {
		rate = c.USD
	}
	if req.ToCurrency == "EUR" {
		rate = c.EUR
	}
	if req.ToCurrency == "RUB" {
		rate = c.RUB
	}

	return &pb.ExchangeRateResponse{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         rate,
	}, nil
}

func startGrpcServer(db *pgxpool.Pool) {
	lis, _ := net.Listen("tcp", ":50051")
	grpcServer := grpc.NewServer()
	pb.RegisterExchangeServiceServer(grpcServer, &ExchangeServer{db: db})
	grpcServer.Serve(lis)
}

func main() {
	connection := fmt.Sprintf("host=postgres port=5432 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	db, err := pgxpool.New(context.Background(), connection)

	if err != nil {
		fmt.Println("ОШИБКА ПРИ ПОДКЛЮЧЕНИИ К БД", err)
		os.Exit(1)
	}

	go startGrpcServer(db)

}
