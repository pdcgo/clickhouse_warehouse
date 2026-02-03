package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

type Tool *cli.Command

func NewTool(
	createForTest CreateForTest,
	migtest Migtest,
	devel Devel,
) Tool {
	return &cli.Command{
		Commands: []*cli.Command{
			{
				Name:   "create_for_test",
				Action: cli.ActionFunc(createForTest),
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "sql",
						Value: true,
					},
				},
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "migration name",
					},
				},
			},
			{
				Name:   "migtest",
				Action: cli.ActionFunc(migtest),
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "command",
					},
					&cli.StringArg{
						Name: "options",
					},
				},
			},
			{
				Name:   "devel",
				Action: cli.ActionFunc(devel),
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "command",
					},
					&cli.StringArg{
						Name: "options",
					},
				},
			},
		},
	}
}

func main() {
	tool, err := InitializeTool()
	if err != nil {
		log.Fatal(err)
	}

	var app *cli.Command = tool
	err = app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
