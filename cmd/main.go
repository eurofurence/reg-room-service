package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/eurofurence/reg-room-service/internal/web/app"
)

func main() {
	// TODO perform initialization

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "foo",
				Aliases: []string{"f"},
				EnvVars: []string{"FOOBAR"},
			},
		},
		Action: func(ctx *cli.Context) error {
			app.New().Run()

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
