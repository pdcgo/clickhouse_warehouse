package main

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/pdcgo/clickhouse_warehouse/replication"
	"github.com/pdcgo/shared/configs"
	"github.com/urfave/cli/v3"
)

func NewReplicationState() (replication.ReplicationState, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		return nil, err
	}
	return replication.NewFirestoreReplicationState(ctx, client, "devel")
}

type GetReplication func(ctx context.Context) (*replication.Replication, error)

func NewReplication(
	cfg *configs.AppConfig,
	state replication.ReplicationState,
) GetReplication {
	return func(ctx context.Context) (*replication.Replication, error) {
		return replication.ConnectReplication(ctx, cfg.Database, state)
	}
}

type AppReplication *cli.Command

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
