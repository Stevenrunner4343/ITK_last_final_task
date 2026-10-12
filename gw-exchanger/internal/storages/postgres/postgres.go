package postgres

import (
	"context"
	"fmt"
	"gw-exchanger/internal/config"
	model "gw-exchanger/internal/storages"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	db *pgxpool.Pool
}

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

func (p *PostgresStorage) GetRates(ctx context.Context) *model.Currency {
	var c model.Currency
	err := p.db.QueryRow(ctx, "SELECT usd, eur, rub FROM currency LIMIT 1").Scan(&c.USD, &c.EUR, &c.RUB)
	if err != nil {
		fmt.Println("ОШИБКА ПРИ получении курсоВ в GetRates ", err)
	}
	return &c
}
