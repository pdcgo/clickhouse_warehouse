package main

import (
	"context"
	"os"

	"github.com/pdcgo/clickhouse_warehouse/replication"
	"github.com/pdcgo/shared/configs"
	"github.com/urfave/cli/v3"
)

type AppReplication *cli.Command

type GetReplication func(ctx context.Context) (*replication.Replication, error)

func NewReplication(
	cfg *configs.AppConfig,
	state replication.ReplicationState,
) GetReplication {
	return func(ctx context.Context) (*replication.Replication, error) {
		return replication.ConnectReplication(ctx, cfg.Database, state)
	}
}

func NewAppReplication(
	start StartFunc,
	backfill BackfillFunc,
) AppReplication {
	return &cli.Command{
		Commands: []*cli.Command{
			{
				Name:        "start",
				Description: "start replication",
				Action:      cli.ActionFunc(start),
			},
			{
				Name:        "backfill",
				Description: "starting backfilling data",
				Action:      cli.ActionFunc(backfill),
			},
		},
	}
}

func main() {
	app, err := InitializeAppReplication()
	if err != nil {
		panic(err)
	}

	var cliApp *cli.Command = app
	err = cliApp.Run(context.Background(), os.Args)
	if err != nil {
		panic(err)
	}

}
