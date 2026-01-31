package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/pressly/goose/v3"
)

func main() {
	db := clickhouse.OpenDB(&clickhouse.Options{
		Addr:     []string{"localhost:9000"}, // 9440 is a secure native TCP port
		Protocol: clickhouse.Native,
		Auth: clickhouse.Auth{
			Username: "user",
			Password: "password",
		},
	})

	defer db.Close()
	// ctx := context.Background()
	err := goose.Create(db, "./test_migrations", "init_test", "sql")
	if err != nil {
		panic(err)
	}

	// err := InitializeVersion(ctx, "test", db)

	// if err != nil {
	// 	panic(err)
	// }
}

type Version struct {
	Version   int64
	IsApplied bool
	Tstamp    string
}

func InitializeVersion(ctx context.Context, mode string, db *sql.DB) error {
	var err error

	slog.Info("creating version")

	// create database

	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS goose_db_version_%s
(
    version_id Int64,
    is_applied UInt8,
    tstamp DateTime DEFAULT now()
)
ENGINE = MergeTree
ORDER BY version_id	
`, mode)

	_, err = db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	rows, err := db.QueryContext(ctx, fmt.Sprintf("select * from goose_db_version_%s", mode))
	if err != nil {
		return err
	}

	var versions []*Version
	for rows.Next() {
		var version Version
		err = rows.Scan(&version.Version, &version.IsApplied, &version.Tstamp)
		if err != nil {
			return err
		}

		versions = append(versions, &version)
		fmt.Println(version)
	}

	if len(versions) == 0 {
		slog.Info("inserting log")
		_, err = db.ExecContext(ctx, fmt.Sprintf("INSERT INTO goose_db_version_%s (version_id, is_applied) VALUES (1, 1)", mode))
		if err != nil {
			return err
		}
	}
	err = goose.UpContext(ctx, db, "./test_migrations")
	return err
}
