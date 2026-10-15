package handlers

import (
	"gw-exchanger/internal/service"

	"context"

	pb "github.com/Stevenrunner4343/proto-exchange/exchange"
	"google.golang.org/grpc"
)

type ExchangeServer struct {
	pb.UnimplementedExchangeServiceServer
	svc *service.ExchangeService
}

func NewExchangeServer(svc *service.ExchangeService) *ExchangeServer {
	return &ExchangeServer{svc: svc}
}

func (s *ExchangeServer) GetExchangeRates(ctx context.Context, empty *pb.Empty) (*pb.ExchangeRatesResponse, error) {
	return s.svc.GetRates(ctx)
}

func (s *ExchangeServer) GetExchangeRateForCurrency(ctx context.Context, req *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	return s.svc.GetRateForCurrency(ctx, req.FromCurrency, req.ToCurrency)
}

func StartGrpcServer(svc *service.ExchangeService) *grpc.Server {
	grpcServer := grpc.NewServer()
	pb.RegisterExchangeServiceServer(grpcServer, NewExchangeServer(svc))
	return grpcServer
}
