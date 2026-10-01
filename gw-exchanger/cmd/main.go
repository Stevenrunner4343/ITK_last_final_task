package main

import (
	"context"
	"fmt"
	"gw-exchanger/internal/storages"
	"net"
	"os"

	pb "github.com/Stevenrunner4343/proto-exchange/exchange"

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
	fmt.Println("сервер запустился")
	grpcServer.Serve(lis)

}

func main() {
	blocker := make(chan struct{}) // без блокера не работает! и клиент обгоняте сервер и поэтому не запускается и
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
	<-blocker
	// можно ли не блокером это сделать а gracefull shutDown
	// stop := make(chan os.Signal, 1)

	// signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	// <-stop
	// bтипо от так

}
