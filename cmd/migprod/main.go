// This is custom goose binary with sqlite3 support only.

package main

import (
	"context"
	"crypto/tls"
	"log"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/pressly/goose/v3"
)

func main() {

	cfg, err := GetProductionClickhouseConfig()
	if err != nil {
		panic(err)
	}

	db := clickhouse.OpenDB(&clickhouse.Options{
		Addr:     []string{cfg.Address}, // 9440 is a secure native TCP port
		Protocol: clickhouse.Native,
		TLS:      &tls.Config{}, // enable secure TLS
		Auth: clickhouse.Auth{
			Username: cfg.Username,
			Password: cfg.Password,
			Database: cfg.Database,
		},
	})

	goose.SetTableName("prod.goose_db_version")
	goose.SetDialect("clickhouse")
	// goose.SetSequential(true)

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v", err)
		}
	}()

	if len(os.Args) < 2 {
		panic("no command")
	}

	command := os.Args[1]
	arguments := []string{}

	if len(os.Args) > 2 {
		arguments = append(arguments, os.Args[2:]...)
	}

	ctx := context.Background()
	dir := "./prod_migrations"
	if err := goose.RunContext(ctx, command, db, dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
