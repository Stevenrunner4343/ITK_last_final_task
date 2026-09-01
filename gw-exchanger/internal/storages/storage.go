package storages

import "context"

type Storage interface {
	GetRates(ctx context.Context) (*Currency, error)
} // не использую, потому что всё отдельно в GRpc функциях
