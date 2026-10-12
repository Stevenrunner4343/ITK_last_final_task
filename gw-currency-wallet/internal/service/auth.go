package service

import (
	"context"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/repository"
	"gw-currency-wallet/internal/storages/model"
	"gw-currency-wallet/pkg/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type AuthService struct {
	repo *repository.WalletRepo
}

func NewAuthService(repo *repository.WalletRepo) *AuthService {
	return &AuthService{repo: repo}
}

func (a *AuthService) MakeNewUser(ctx context.Context, acc model.Account) (string, error) {
	userId, err := a.repo.MakeNewUser(ctx, acc)
	if err != nil {
		logger.Error("Ошибка при Создании нового Юзера", zap.Error(err))
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(config.JWT)
	return tokenString, nil

}
func (a *AuthService) CheckUser(ctx context.Context, acc model.Account) error {
	err := a.repo.CheckUser(ctx, acc)
	if err != nil {
		if err == repository.ErrNotFound {
			logger.Error("Юзер не Найден", zap.Error(err))
			return repository.ErrNotFound
		} else if err == repository.ErrWrongPassword {
			logger.Error("Не Верный пароль", zap.Error(err))
			return repository.ErrWrongPassword
		} else {
			logger.Error("Ошьика при записи в БД ", zap.Error(err))
			return err
		}
	} else {
		return err
	}

}
