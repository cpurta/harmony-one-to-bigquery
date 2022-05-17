package command

import (
	"time"

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
				EnvVars:     []string{"NODE_URL"},
				Usage:       "the url of the node used to pull historical data from",
				Destination: &backfillRunner.NodeURL,
				Value:       "https://api.s0.t.hmny.io",
			},
			&cli.StringFlag{
				Name:        "gcp-project-id",
				EnvVars:     []string{"GCP_PROJECT_ID"},
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
				EnvVars:     []string{"GCP_DATASET_ID"},
				Usage:       "the dataset id used in GCP to store blockchain data in BigQuery",
				Destination: &backfillRunner.DatasetID,
				Value:       "crypto_harmony",
			},
			&cli.StringFlag{
				Name:        "gcp-blocks-table-id",
				EnvVars:     []string{"GCP_BLOCKS_TABLE_ID"},
				Usage:       "the blocks table id used in GCP to store blockchain data in BigQuery",
				Destination: &backfillRunner.BlocksTableID,
				Value:       "blocks",
			},
			&cli.StringFlag{
				Name:        "gcp-txns-table-id",
				EnvVars:     []string{"GCP_TXNS_TABLE_ID"},
				Usage:       "the transactions table id used in GCP to store blockchain data in BigQuery",
				Destination: &backfillRunner.TxnsTableID,
				Value:       "transactions",
			},
			&cli.BoolFlag{
				Name:        "logging-production",
				Usage:       "determines the logging level and format to seperate development and production environment",
				Destination: &backfillRunner.LoggingProduction,
			},
			&cli.StringFlag{
				Name:        "log-level",
				EnvVars:     []string{"LOG_LEVEL"},
				Usage:       "the logging level used in outputting logs",
				Destination: &backfillRunner.LogLevel,
				Value:       "info",
			},
			&cli.Int64Flag{
				Name:        "start-block",
				EnvVars:     []string{"START_BLOCK"},
				Usage:       "start backfill from a specific block, if the most start block exceeds the most recent synced block it will fail",
				Destination: &backfillRunner.StartBlock,
				Value:       0,
			},
			&cli.IntFlag{
				Name:        "concurrency",
				EnvVars:     []string{"CONCURRENCY"},
				Usage:       "the number concurrent go routines pulling Harmony One blockchain data",
				Destination: &backfillRunner.Concurrency,
				Value:       1,
			},
			&cli.IntFlag{
				Name:        "max-retries",
				EnvVars:     []string{"MAX_RETRIES"},
				Usage:       "the maximum number to times a block or transaction will be attempted to be inserted into their respective tables",
				Destination: &backfillRunner.MaxRetries,
				Value:       10,
			},
			&cli.DurationFlag{
				Name:        "wait-time",
				EnvVars:     []string{"WAIT_TIME"},
				Usage:       "the amount of time to wait before querying for the most recent block",
				Destination: &backfillRunner.WaitTime,
				Value:       time.Duration(time.Second * 30),
			},
		},
		Action: backfillRunner.Run,
	}

	return cmd
}
