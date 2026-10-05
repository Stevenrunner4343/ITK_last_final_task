package clickhouse

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type AggRow struct {
	Period    string
	Operation string
	Cnt       uint64
	AvgLat    float64
}

func Aggregation(ctx context.Context, conn driver.Conn) {
	periods := []struct {
		name  string
		query string
	}{
		{
			"1 MINUTE",
			`SELECT toString(toStartOfMinute(created_at)) AS period,
			        operation,
			        count() AS cnt,
			        round(avg(latency_ms), 2) AS avg_lat
			 FROM events FINAL
			 GROUP BY period, operation
			 ORDER BY period DESC
			 LIMIT 10`,
		},
		{
			"5 MINUTES",
			`SELECT toString(toStartOfInterval(created_at, INTERVAL 5 MINUTE)) AS period,
			        operation,
			        count() AS cnt,
			        round(avg(latency_ms), 2) AS avg_lat
			 FROM events FINAL
			 GROUP BY period, operation
			 ORDER BY period DESC
			 LIMIT 10`,
		},
		{
			"1 HOUR",
			`SELECT toString(toStartOfHour(created_at)) AS period,
			        operation,
			        count() AS cnt,
			        round(avg(latency_ms), 2) AS avg_lat
			 FROM events FINAL
			 GROUP BY period, operation
			 ORDER BY period DESC
			 LIMIT 10`,
		},
		{
			"1 DAY",
			`SELECT toString(toStartOfDay(created_at)) AS period,
			        operation,
			        count() AS cnt,
			        round(avg(latency_ms), 2) AS avg_lat
			 FROM events FINAL
			 GROUP BY period, operation
			 ORDER BY period DESC
			 LIMIT 10`,
		},
		{
			"1 WEEK",
			`SELECT toString(toStartOfWeek(created_at)) AS period,
			        operation,
			        count() AS cnt,
			        round(avg(latency_ms), 2) AS avg_lat
			 FROM events FINAL
			 GROUP BY period, operation
			 ORDER BY period DESC
			 LIMIT 10`,
		},
	}

	for _, p := range periods {
		fmt.Printf("\n=== %s ===\n", p.name)

		rows, err := conn.Query(ctx, p.query)
		if err != nil {
			fmt.Println("query error:", err)
			continue
		}

		fmt.Printf("%-22s %-12s %8s %10s\n", "period", "operation", "cnt", "avg_lat")
		for rows.Next() {
			var r AggRow
			if err := rows.Scan(&r.Period, &r.Operation, &r.Cnt, &r.AvgLat); err != nil {
				fmt.Println("scan error:", err)
				continue
			}
			fmt.Printf("%-22s %-12s %8d %10.2f\n", r.Period, r.Operation, r.Cnt, r.AvgLat)
		}
		_ = rows.Close()
	}
}
