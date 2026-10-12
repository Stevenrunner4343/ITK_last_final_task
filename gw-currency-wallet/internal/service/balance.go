package service

import (
	"context"
	"errors"
	"gw-currency-wallet/internal/repository"
	"gw-currency-wallet/pkg/logger"

	"go.uber.org/zap"
)

type BalanceService struct {
	repo *repository.WalletRepo
}

var (
	ErrUserNotFound = errors.New("Юзер не найден")
)

func NewBalanceService(repo *repository.WalletRepo) *BalanceService {
	return &BalanceService{repo: repo}
}

func (s *BalanceService) GetBalance(ctx context.Context, Userid int) (int, error) {

	balance, err := s.repo.FindUser(ctx, Userid)

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return 0, ErrUserNotFound
		}
		logger.Error("Ошибка при запросе к БД на Получение Баланса", zap.Error(err))
		return 0, err
	} else {
		return balance, nil
	}
}
