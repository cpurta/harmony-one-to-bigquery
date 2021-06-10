package main

import (
	"fmt"
	"os"

	command "github.com/cpurta/harmony-one-to-bigquery/internal/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "hmy-bq-import",
		Usage: "command line tool to import Harmony One blockchain data into GCP BigQuery",
		Commands: []*cli.Command{
			command.BackfillCommand(),
		},
		Version: "v0.0.1",
		Authors: []*cli.Author{
			{
				Name:  "Chris Purta",
				Email: "cpurta@gmail.com",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("error running program:", err.Error())
		os.Exit(1)
	}
}
