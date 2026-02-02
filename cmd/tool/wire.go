//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/urfave/cli/v3"
)

func InitializeTool() (Tool, error) {
	wire.Build(
		NewCreateFortest,
		NewMigtest,
		NewTool,
	)
	return &cli.Command{}, nil
}
