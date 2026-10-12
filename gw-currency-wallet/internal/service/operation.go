package service

import (
	"context"
	"gw-currency-wallet/internal/producer"
	"gw-currency-wallet/internal/repository"
	"gw-currency-wallet/internal/storages/model"
	"gw-currency-wallet/pkg/logger"
	"math/rand/v2"
	"time"

	"go.uber.org/zap"
)

type OperationService struct {
	repo *repository.WalletRepo
	prod *producer.Producer
}

func NewOperationService(repo *repository.WalletRepo, producer *producer.Producer) *OperationService {
	return &OperationService{
		repo: repo,
		prod: producer,
	}
}

func (p *OperationService) Process(ctx context.Context, d model.Data) error {

	if d.Amount >= 30000 {
		err := p.prod.Send(ctx, d, producer.TopicBigMoney)
		if err != nil {
			logger.Error("Ошибка при отправки Брокером сообщения о Боольшой сумме в Notification", zap.Error(err), zap.Int("user_id", d.Id))
			return err
		}
	}

	err := p.repo.SendMoney(ctx, d.Id, d.Operation, d.Amount)
	if err != nil {
		logger.Error("Ошибка при ЗАписи в БД СуММЫ МЕНЬШЕ 30 000", zap.Error(err), zap.Int("user_id", d.Id))
		return err
	}

	var e model.Event
	e.ID = rand.Uint64()
	e.Operation = d.Operation
	e.UserID = d.Id
	e.Amount = d.Amount
	e.CreatedAt = time.Now()

	err = p.prod.Send(ctx, e, producer.TopicAnalytics)
	if err != nil {
		logger.Error("Ошибка при отправки Брокером сообщения в Аналитику", zap.Error(err), zap.Int("user_id", d.Id))
		return err
	}
	return nil
}
