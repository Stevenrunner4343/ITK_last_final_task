package service

import (
	"context"

	"gw-exchanger/internal/repository"

	pb "github.com/Stevenrunner4343/proto-exchange/exchange"
)

type ExchangeService struct {
	repo *repository.ExchangeRepo
}

func NewExchangeService(repo *repository.ExchangeRepo) *ExchangeService {
	return &ExchangeService{repo: repo}
}

func (s *ExchangeService) GetRates(ctx context.Context) (*pb.ExchangeRatesResponse, error) {
	c, err := s.repo.GetRates(ctx)
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

func (s *ExchangeService) GetRateForCurrency(ctx context.Context, from, to string) (*pb.ExchangeRateResponse, error) {
	c, err := s.repo.GetRates(ctx)
	if err != nil {
		return nil, err
	}

	var rate int32
	switch to {
	case "USD":
		rate = c.USD
	case "EUR":
		rate = c.EUR
	case "RUB":
		rate = c.RUB
	}

	return &pb.ExchangeRateResponse{
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         rate,
	}, nil
}
