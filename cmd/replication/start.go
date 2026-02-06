package main

import (
	"context"
	"log/slog"

	"github.com/pdcgo/clickhouse_warehouse/replication"
	"github.com/pdcgo/shared/pkg/debugtool"
	"github.com/urfave/cli/v3"
)

type StartFunc cli.ActionFunc

func NewStartFunc(
	getReplication GetReplication,
) StartFunc {
	return func(ctx context.Context, c *cli.Command) error {

		rep, err := getReplication(ctx)
		if err != nil {
			return err
		}

		defer rep.Close(ctx)

		lsn, err := rep.GetLsn(ctx)
		if err != nil {
			return err
		}

		slog.Info("current lsn", "lsn", lsn.String())

		err = rep.StreamStart(
			ctx,
			"devel_slot",
			"warehouse_publication",
			func(ctx context.Context, event *replication.ReplicationEvent) error {
				switch event.SourceMetadata.Table {
				case "orders":
					debugtool.LogJson(event)
				}

				return nil
			},
		)
		return err
	}
}
