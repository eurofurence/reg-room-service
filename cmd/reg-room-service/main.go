package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/eurofurence/reg-room-service/internal/web/app"
)

var (
	configFilePath  string
	migrateDatabase bool
)

func main() {
	application := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				EnvVars:     []string{"CONFIG_FILE"},
				Required:    true,
				Destination: &configFilePath,
			},
			&cli.BoolFlag{
				Name:        "migrate-database",
				Aliases:     []string{"m"},
				Destination: &migrateDatabase,
			},
		},
		Action: func(ctx *cli.Context) error {
			return app.New(
				app.NewParams(
					configFilePath,
					migrateDatabase,
				)).Run()
		},
	}

	if err := application.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
