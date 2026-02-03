package main

import (
	"context"
	"errors"

	"github.com/pdcgo/clickhouse_warehouse/database"
	"github.com/pressly/goose/v3"
	"github.com/urfave/cli/v3"
)

type CreateForTest cli.ActionFunc

func NewCreateFortest() CreateForTest {
	return func(ctx context.Context, c *cli.Command) error {
		db := database.NewLocalDatabase()
		isSql := c.Bool("sql")

		fname := c.Arguments[0].Get().(string)
		if fname == "" {
			return errors.New("migration name required")
		}

		var migrationType string
		if isSql {
			migrationType = "sql"
		} else {
			migrationType = "go"
		}

		return goose.Create(db, "./test_migrations", fname, migrationType)

	}
}
