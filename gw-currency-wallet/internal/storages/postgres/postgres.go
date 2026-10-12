package postgres

import (
	"context"
	"fmt"
	"gw-currency-wallet/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnection(ctx context.Context) (*pgxpool.Pool, error) {
	conn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
	)

	db, err := pgxpool.New(ctx, conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
