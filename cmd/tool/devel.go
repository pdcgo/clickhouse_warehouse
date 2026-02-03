package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/pdcgo/clickhouse_warehouse/database"
	"github.com/pressly/goose/v3"
	"github.com/urfave/cli/v3"
)

type Devel cli.ActionFunc

func NewDevel() Devel {

	return func(ctx context.Context, c *cli.Command) error {

		db := database.NewLocalDatabase()

		command := c.Arguments[0].Get().(string)
		args := []string{}
		for _, arg := range c.Arguments[1:] {
			args = append(args, arg.Get().(string))
		}

		// running initialization
		versionTableName := "dbversion_devel"
		err := InitializeVersion(ctx, versionTableName, db)
		if err != nil {
			return err
		}
		goose.SetTableName(versionTableName)
		goose.SetDialect("clickhouse")

		// running migration
		dir := "./devel_migrations"

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
