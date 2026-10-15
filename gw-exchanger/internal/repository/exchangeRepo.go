package repository

import (
	"context"

	"gw-exchanger/internal/storages/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ExchangeRepo struct {
	db *pgxpool.Pool
}

func NewExchangeRepo(db *pgxpool.Pool) *ExchangeRepo {
	return &ExchangeRepo{db: db}
}

func (r *ExchangeRepo) GetRates(ctx context.Context) (model.Currency, error) {
	var c model.Currency
	err := r.db.QueryRow(ctx, "SELECT usd, eur, rub FROM currency").
		Scan(&c.USD, &c.EUR, &c.RUB)
	return c, err
}
