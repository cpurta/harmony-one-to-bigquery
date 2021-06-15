package runner

import (
	"context"
	"errors"
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
	Concurrency          int64
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

	if runner.Concurrency <= 0 {
		return errors.New("must provide positive concurrency value")
	}

	return runner.backfillFromLatest(ctx)
}

func (runner *BackfillRunner) backfillFromLatest(ctx context.Context) error {
	var (
		header             *model.Header
		currentBlock       int64
		backfillUpTo       = int64(0)
		totalBlocks        int64
		blocksPerPartition int64
		err                error
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

	totalBlocks = backfillUpTo - currentBlock

	blocksPerPartition = totalBlocks / runner.Concurrency

	for i := int64(1); i <= runner.Concurrency; i++ {
		if i == runner.Concurrency {
			runner.backfillPartition(ctx, currentBlock+(i-1)*blocksPerPartition, backfillUpTo, i)
		}

		runner.backfillPartition(ctx, currentBlock+(i-1)*blocksPerPartition, currentBlock+i*blocksPerPartition, i)
	}

	return nil
}

func (runner *BackfillRunner) backfillPartition(ctx context.Context, startBlock, endBlock, partition int64) error {
	partitionLogger := runner.logger.With(zap.Int64("partition", partition), zap.Int64("start_block", startBlock), zap.Int64("end_block", endBlock))

	partitionLogger.Info("starting partition backfill")

	for currentBlock := startBlock; currentBlock <= endBlock; currentBlock++ {
		var (
			block *model.Block
			err   error
		)

		if block, err = runner.harmonyClient.GetBlockByNumber(currentBlock); err != nil {
			runner.logger.Error("unable get block from harmony blockchain client", zap.Int64("block_number", currentBlock), zap.Error(err))
			continue
		}

		if runner.DryRun {
			runner.logger.Info("received block", zap.Int64("block_number", currentBlock))
			continue
		}

		if err = runner.bigQueryClient.InsertBlock(block, ctx); err != nil {
			runner.logger.Error("unable to insert block into BigQuery", zap.Int64("block_number", currentBlock), zap.Error(err))
		}

		if err = runner.bigQueryClient.InsertTransactions(block.Transactions, ctx); err != nil {
			runner.logger.Error("unable to insert block transactions into BigQuery", zap.Int64("block_number", currentBlock), zap.Error(err))
		}
	}

	return nil
}
