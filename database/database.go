package database

import (
	"database/sql"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func NewLocalDatabase() *sql.DB {
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

func NewLocalDatabaseHttp() (driver.Conn, error) {

	return clickhouse.Open(&clickhouse.Options{
		Protocol: clickhouse.HTTP,
		Addr:     []string{"localhost:8123"}, // 9440 is a secure native TCP port
		Auth: clickhouse.Auth{
			Username: "user",
			Password: "password",
		},
	})

}
