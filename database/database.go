package database

import (
	"database/sql"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func NewTestDatabase() *sql.DB {
	db := clickhouse.OpenDB(&clickhouse.Options{
		Addr:     []string{"localhost:9000"}, // 9440 is a secure native TCP port
		Protocol: clickhouse.Native,
		Auth: clickhouse.Auth{
			Username: "user",
			Password: "password",
		},
	})
	return db
}
