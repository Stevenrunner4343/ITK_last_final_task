package storages

import "context"

type Storage interface {
	GetRates(ctx context.Context) (*Currency, error)
}
