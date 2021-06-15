package runner

import (
	"context"
	"net/http"

	"github.com/cpurta/harmony-one-to-bigquery/internal/clients/bigquery"
	bq "github.com/cpurta/harmony-one-to-bigquery/internal/clients/bigquery/client"
	"github.com/cpurta/harmony-one-to-bigquery/internal/clients/harmony"
	"github.com/cpurta/harmony-one-to-bigquery/internal/clients/harmony/client"
	"github.com/cpurta/harmony-one-to-bigquery/internal/model"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type BackfillRunner struct {
	NodeURL              string
	GoogleCloudProjectID string
	DryRun               bool
	DatasetID            string
	BlocksTableID        string
	TxnsTableID          string
	harmonyClient        harmony.HarmonyClient
	bigQueryClient       bigquery.BigQueryClient
	logger               *zap.Logger
}

func (runner *BackfillRunner) Run(cli *cli.Context) error {
	var (
		ctx = context.Background()
		err error
	)

	if runner.logger, err = zap.NewProduction(); err != nil {
		return err
	}

	runner.harmonyClient = client.NewHarmonyOneClient(runner.NodeURL, http.DefaultClient, runner.logger.Named("harmony_client"))

	if runner.bigQueryClient, err = bq.NewBigQueryClient(ctx, runner.GoogleCloudProjectID, runner.DatasetID, runner.BlocksTableID, runner.TxnsTableID); err != nil {
		return err
	}

	return runner.backfillFromLatest(ctx)
}

func (runner *BackfillRunner) backfillFromLatest(ctx context.Context) error {
	var (
		header       *model.Header
		currentBlock int64
		backfillUpTo = int64(0)
		err          error
	)

	if header, err = runner.harmonyClient.GetLatestHeader(); err != nil {
		runner.logger.Error("unable to get the most recent block header from the harmony blockchain client", zap.Error(err))
		return err
	}

	runner.logger.Info("received current block header on hmy blockchain", zap.Int64("hmy_block_number", header.BlockNumber))

	if currentBlock, err = runner.bigQueryClient.GetMostRecentBlockNumber(ctx); err != nil {
		runner.logger.Error("unable to get the most recent block number stored in BigQuery", zap.Error(err))
		return err
	}

	runner.logger.Info("received most recent block number stored in BigQuery", zap.Int64("bq_block_number", currentBlock))

	backfillUpTo = header.BlockNumber

	for ; currentBlock <= backfillUpTo; currentBlock++ {
		var (
			block             *model.Block
			blockNumberLogger = runner.logger.With(zap.Int64("block_number", currentBlock))
		)

		if block, err = runner.harmonyClient.GetBlockByNumber(currentBlock); err != nil {
			blockNumberLogger.Error("unable get block from harmony blockchain client", zap.Error(err))
			continue
		}

		if currentBlock%1000 == 0 {
			blockNumberLogger.Info("processed another 1000 blocks")
		}

		if runner.DryRun {
			blockNumberLogger.Info("received block")
			continue
		}

		if err = runner.bigQueryClient.InsertBlock(block, ctx); err != nil {
			blockNumberLogger.Error("unable to insert block into BigQuery", zap.Error(err))
		}

		if err = runner.bigQueryClient.InsertTransactions(block.Transactions, ctx); err != nil {
			blockNumberLogger.Error("unable to insert block transactions into BigQuery", zap.Error(err))
		}
	}

	return nil
}
