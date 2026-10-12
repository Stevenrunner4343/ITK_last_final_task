package handlers

import (
	"context"
	"fmt"
	model "gw-exchanger/internal/storages"
	"net"

	pb "github.com/Stevenrunner4343/proto-exchange/exchange"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type ExchangeServer struct {
	pb.UnimplementedExchangeServiceServer
	db *pgxpool.Pool
}

func (s *ExchangeServer) GetExchangeRates(ctx context.Context, empty *pb.Empty) (*pb.ExchangeRatesResponse, error) {
	var c model.Currency

	err := s.db.QueryRow(ctx, "SELECT * FROM currency").
		Scan(&c.USD, &c.EUR, &c.RUB)
	if err != nil {
		return nil, err
	}

	return &pb.ExchangeRatesResponse{
		Rates: map[string]int32{
			"USD": c.USD,
			"EUR": c.EUR,
			"RUB": c.RUB,
		},
	}, nil
}

func (s *ExchangeServer) GetExchangeRateForCurrency(ctx context.Context, req *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	var c model.Currency

	err := s.db.QueryRow(ctx, "SELECT * FROM currency").
		Scan(&c.USD, &c.EUR, &c.RUB)
	if err != nil {
		return nil, err
	}

	var rate int32

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

func StartGrpcServer(db *pgxpool.Pool) {
	lis, _ := net.Listen("tcp", ":50051")
	grpcServer := grpc.NewServer()
	pb.RegisterExchangeServiceServer(grpcServer, &ExchangeServer{db: db})
	fmt.Println("сервер запустился")
	grpcServer.Serve(lis)

}
