package service

import (
	"context"
	"encoding/json"

	"gw-notification/internal/repository"
	"gw-notification/internal/storages/model"
	"gw-notification/pkg/logger"

	"go.uber.org/zap"
)

type NotificationService struct {
	repo *repository.NotificationRepo
}

func NewNotificationService(repo *repository.NotificationRepo) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) ProcessMessage(ctx context.Context, raw []byte) error {
	var d model.Data
	err := json.Unmarshal(raw, &d)

	err = s.repo.SendToMongo(ctx, d)
	if err != nil {
		logger.Error("ошибка при Вызове функции отправки в МОнго", zap.Error(err))
		return err
	}

	return nil
}
