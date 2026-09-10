package storages

import (
	"context"
)

type Storage interface {
	SendToMongo(ctx context.Context, name string, age int)
}
