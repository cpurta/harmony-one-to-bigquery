package command

import (
	"github.com/cpurta/harmony-one-to-bigquery/internal/runner"
	"github.com/urfave/cli/v2"
)

func BackfillCommand() *cli.Command {
	var (
		backfillRunner = &runner.BackfillRunner{}
	)

	cmd := &cli.Command{
		Name:  "backfill",
		Usage: "pull historical block data from Harmony One blockchain and insert into GCP BigQuery",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "node-url",
				Usage:       "the url of the node used to pull historical data from",
				Destination: &backfillRunner.NodeURL,
				Value:       "https://api.s0.t.hmny.io",
			},
			&cli.StringFlag{
				Name:        "gcp-project-id",
				Usage:       "the project id used in GCP to store blockchain data in BigQuery",
				Destination: &backfillRunner.GoogleCloudProjectID,
			},
			&cli.BoolFlag{
				Name:        "dry-run",
				Usage:       "pull historical blockchain data but do not attempt to insert it into BigQuery",
				Destination: &backfillRunner.DryRun,
			},
			&cli.StringFlag{
				Name:        "gcp-dataset-id",
				Usage:       "the dataset id used in GCP to store blockchain data in BigQuery",
				Destination: &backfillRunner.DatasetID,
				Value:       "crypto_harmony",
			},
			&cli.StringFlag{
				Name:        "gcp-blocks-table-id",
				Usage:       "the blocks table id used in GCP to store blockchain data in BigQuery",
				Destination: &backfillRunner.BlocksTableID,
				Value:       "blocks",
			},
			&cli.StringFlag{
				Name:        "gcp-txns-table-id",
				Usage:       "the transactions table id used in GCP to store blockchain data in BigQuery",
				Destination: &backfillRunner.TxnsTableID,
				Value:       "transactions",
			},
		},
		Action: backfillRunner.Run,
	}

	return cmd
}
