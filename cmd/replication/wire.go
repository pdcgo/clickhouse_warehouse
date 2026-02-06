//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/pdcgo/clickhouse_warehouse/replication"
	"github.com/pdcgo/shared/configs"
	"github.com/urfave/cli/v3"
)

func InitializeAppReplication() (AppReplication, error) {
	wire.Build(
		configs.NewProductionConfig,
		replication.NewMemoryReplicationState,
		NewReplication,
		NewBackfillFunc,
		NewStartFunc,
		NewAppReplication,
	)
	return &cli.Command{}, nil
}
