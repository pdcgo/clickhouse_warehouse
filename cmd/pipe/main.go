package main

import (
	"crypto/tls"
	"database/sql"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func NewLocalDatabase() *sql.DB {
	return clickhouse.OpenDB(&clickhouse.Options{
		Addr:     []string{"localhost:9000"}, // 9440 is a secure native TCP port
		Protocol: clickhouse.Native,
		TLS:      &tls.Config{}, // enable secure TLS
		Auth: clickhouse.Auth{
			Username: "user",
			Password: "password",
			Database: "default",
		},
	})
}

func main() {
	db := NewLocalDatabase()
	defer db.Close()

}

type Pipeline struct {
	Down string
	Up   string
}
