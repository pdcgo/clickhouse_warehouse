//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/pdcgo/shared/configs"
	"github.com/urfave/cli/v3"
)

func InitializeAppReplication() (AppReplication, error) {
	wire.Build(
		configs.NewProductionConfig,
		NewReplicationState,
		NewReplication,
		NewBackfillFunc,
		NewStartFunc,
		NewAppReplication,
	)
	return &cli.Command{}, nil
}
