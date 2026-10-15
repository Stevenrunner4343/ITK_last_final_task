package service

import (
	"context"
	"gw-currency-wallet/internal/client"
	"gw-currency-wallet/pkg/logger"

	pb "github.com/Stevenrunner4343/proto-exchange/exchange"
	"go.uber.org/zap"
)

type RatesService struct {
	conn *client.GrpcClient
}

func NewRatesService(client *client.GrpcClient) *RatesService {
	return &RatesService{conn: client}
}
func (r *RatesService) Rates(ctx context.Context) (map[string]int32, error) {

	Rates, err := r.conn.GetExchangeRates(ctx, &pb.Empty{})
	if err != nil {
		logger.Error("Ошибка при Подключении Данных Rates", zap.Error(err))
		return nil, err
	}
	return Rates.Rates, nil

}
func (r *RatesService) Rate(ctx context.Context, from string, to string) (int32, error) {
	Rate, err := r.conn.GetExchangeRateForCurrency(ctx, &pb.CurrencyRequest{
		FromCurrency: from,
		ToCurrency:   to,
	})

	if err != nil {
		logger.Error("Ошибка при Получении Данных Rate", zap.Error(err))
		return 0, err
	}
	return Rate.Rate, nil

}
