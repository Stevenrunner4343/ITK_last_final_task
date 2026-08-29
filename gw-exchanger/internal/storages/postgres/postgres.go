package postgres

import (
	"context"
	"fmt"
	"gw-exchanger/internal/storages"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (p *PostgresStorage) GetRates(ctx context.Context) *storages.Currency {
	var c storages.Currency
	err := p.db.QueryRow(ctx, "SELECT usd, eur, rub FROM currency LIMIT 1").Scan(&c.USD, &c.EUR, &c.RUB)
	if err != nil {
		fmt.Println("ОШИБКА ПРИ получении курсоВ в GetRates ", err)
	}
	return &c
}
