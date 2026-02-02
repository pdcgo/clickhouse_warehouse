package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/pdcgo/clickhouse_warehouse/database"
	"github.com/pressly/goose/v3"
	"github.com/urfave/cli/v3"
)

type Migtest cli.ActionFunc

func NewMigtest() Migtest {

	return func(ctx context.Context, c *cli.Command) error {

		db := database.NewTestDatabase()

		command := c.Arguments[0].Get().(string)
		args := []string{}
		for _, arg := range c.Arguments[1:] {
			args = append(args, arg.Get().(string))
		}

		// running initialization
		versionTableName := "dbversion_test"
		err := InitializeVersion(ctx, versionTableName, db)
		if err != nil {
			return err
		}
		goose.SetTableName(versionTableName)
		goose.SetDialect("clickhouse")

		// running migration
		dir := "./test_migrations"

		switch command {
		case "create":
			if len(c.Arguments) < 1 {
				return errors.New("fname needed")
			}

			fname := c.Arguments[1].Get().(string)
			if err := goose.Create(db, dir, fname, "sql"); err != nil {
				return err
			}
		case "up":
			if err := goose.UpContext(ctx, db, dir); err != nil {
				return err
			}
		case "up-by-one":
			if err := goose.UpByOneContext(ctx, db, dir); err != nil {
				return err
			}
		case "up-to":
			if len(args) == 0 {
				return fmt.Errorf("up-to must be of form: goose [OPTIONS] DRIVER DBSTRING up-to VERSION")
			}

			version, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("version must be a number (got '%s')", args[0])
			}
			if err := goose.UpToContext(ctx, db, dir, version); err != nil {
				return err
			}
		case "down":
			if err := goose.DownContext(ctx, db, dir); err != nil {
				return err
			}
		case "down-to":
			if len(args) == 0 {
				return fmt.Errorf("down-to must be of form: goose [OPTIONS] DRIVER DBSTRING down-to VERSION")
			}

			version, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("version must be a number (got '%s')", args[0])
			}
			if err := goose.DownToContext(ctx, db, dir, version); err != nil {
				return err
			}
		case "fix":
			if err := goose.Fix(dir); err != nil {
				return err
			}
		case "redo":
			if err := goose.RedoContext(ctx, db, dir); err != nil {
				return err
			}
		case "reset":
			if err := goose.ResetContext(ctx, db, dir); err != nil {
				return err
			}
		case "status":
			if err := goose.StatusContext(ctx, db, dir); err != nil {
				return err
			}
		case "version":
			if err := goose.VersionContext(ctx, db, dir); err != nil {
				return err
			}
		}

		return nil
	}
}

type Version struct {
	Version   int64
	IsApplied bool
	Tstamp    string
}

func InitializeVersion(ctx context.Context, versionTableName string, db *sql.DB) error {
	var err error

	slog.Info("creating version")

	// create database

	query := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s
(
    version_id Int64,
    is_applied UInt8,
    tstamp DateTime DEFAULT now()
)
ENGINE = MergeTree
ORDER BY version_id	
`, versionTableName)

	_, err = db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	rows, err := db.QueryContext(ctx, fmt.Sprintf("select * from %s", versionTableName))
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
		_, err = db.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s (version_id, is_applied) VALUES (1, 1)", versionTableName))
		if err != nil {
			return err
		}
	}

	return err
}
