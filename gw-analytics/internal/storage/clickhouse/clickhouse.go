package clickhouse

import (
	"gw-analytics/internal/config"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func NewConnection() (driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: config.CLICKHOUSE_ADDR,

		Auth: clickhouse.Auth{
			Database: config.CLICKHOUSE_DB,
			Username: config.CLICKHOUSE_USER,
			Password: config.CLICKHOUSE_PASSWORD,
		},
	})
	if err != nil {
		return nil, err
	}
	return conn, nil
}
