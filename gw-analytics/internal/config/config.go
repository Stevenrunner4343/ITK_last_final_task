package config

import "os"

var (
	CLICKHOUSE_DB       = os.Getenv("CLICKHOUSE_DB")
	CLICKHOUSE_USER     = os.Getenv("CLICKHOUSE_USER")
	CLICKHOUSE_PASSWORD = os.Getenv("CLICKHOUSE_PASSWORD")

	CLICKHOUSE_ADDR = []string{os.Getenv("CLICKHOUSE_ADDR")}
)
