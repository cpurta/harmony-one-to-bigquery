package runner

import (
	"context"
	"errors"
	"net/http"
	"time"

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
	FromBlock            int64
	ToBlock              int64
}

func (runner *BackfillRunner) Run(cli *cli.Context) error {
	var (
		harmonyClient  harmony.HarmonyClient
		bigQueryClient bigquery.BigQueryClient
		logger         *zap.Logger
		ctx            = context.Background()
		err            error
	)

	if logger, err = zap.NewProduction(); err != nil {
		return err
	}

	harmonyClient = client.NewHarmonyOneClient(runner.NodeURL, http.DefaultClient)

	if bigQueryClient, err = bq.NewBigQueryClient(ctx, runner.GoogleCloudProjectID, runner.DatasetID, runner.BlocksTableID, runner.TxnsTableID); err != nil {
		return err
	}

	// sanity check to catch from block specified but not to block (or visa versa)
	if runner.FromBlock != int64(-1) && runner.ToBlock == int64(-1) {
		return errors.New("must specify positive value --to-block value")
	}

	if runner.ToBlock != int64(-1) && runner.FromBlock == int64(-1) {
		return errors.New("must specify positive value --from-block value")
	}

	// check if specified to backfill a portion otherwise default to latest
	if runner.FromBlock != int64(-1) && runner.ToBlock != int64(-1) {
		return runner.backfillPortion(ctx, harmonyClient, bigQueryClient, logger)
	}

	return runner.backfillFromLatest(ctx, harmonyClient, bigQueryClient, logger)
}

func (runner *BackfillRunner) backfillPortion(ctx context.Context, harmonyClient harmony.HarmonyClient, bigQueryClient bigquery.BigQueryClient, logger *zap.Logger) error {
	var (
		currentBlock = runner.FromBlock
		totalBlocks  = runner.ToBlock - runner.FromBlock
		start        = time.Now()
		err          error
	)

	for ; currentBlock <= runner.ToBlock; currentBlock++ {
		var block *model.Block

		if block, err = harmonyClient.GetBlockByNumber(currentBlock); err != nil {
			logger.Error("unable get block from harmony blockchain client", zap.Int64("block_number", currentBlock), zap.Error(err))
			continue
		}

		if runner.DryRun {
			logger.Info("received block", zap.Int64("block_number", currentBlock))
			continue
		}

		if err = bigQueryClient.InsertBlock(block, ctx); err != nil {
			logger.Error("unable to insert block into BigQuery", zap.Int64("block_number", currentBlock), zap.Error(err))
		}

		if err = bigQueryClient.InsertTransactions(block.Transactions, ctx); err != nil {
			logger.Error("unable to insert block transactions into BigQuery", zap.Int64("block_number", currentBlock), zap.Error(err))
		}
	}

	logger.Info("complete portion backfill", zap.Int64("total_blocks_backfilled", totalBlocks), zap.Duration("time_elapsed", time.Since(start)))

	return nil
}

func (runner *BackfillRunner) backfillFromLatest(ctx context.Context, harmonyClient harmony.HarmonyClient, bigQueryClient bigquery.BigQueryClient, logger *zap.Logger) error {
	var (
		header       *model.Header
		currentBlock int64
		backfillUpTo = int64(0)
		err          error
	)

	if header, err = harmonyClient.GetLatestHeader(); err != nil {
		logger.Error("unable to get the most recent block header from the harmony blockchain client", zap.Error(err))
		return err
	}

	logger.Info("received current block header on hmy blockchain", zap.Int64("hmy_block_number", header.BlockNumber))

	if currentBlock, err = bigQueryClient.GetMostRecentBlockNumber(ctx); err != nil {
		logger.Error("unable to get the most recent block number stored in BigQuery", zap.Error(err))
		return err
	}

	logger.Info("received most recent block number stored in BigQuery", zap.Int64("bq_block_number", currentBlock))

	backfillUpTo = header.BlockNumber

	for ; currentBlock <= backfillUpTo; currentBlock++ {
		var block *model.Block

		if block, err = harmonyClient.GetBlockByNumber(currentBlock); err != nil {
			logger.Error("unable get block from harmony blockchain client", zap.Int64("block_number", currentBlock), zap.Error(err))
			continue
		}

		if runner.DryRun {
			logger.Info("received block", zap.Int64("block_number", currentBlock))
			continue
		}

		if err = bigQueryClient.InsertBlock(block, ctx); err != nil {
			logger.Error("unable to insert block into BigQuery", zap.Int64("block_number", currentBlock), zap.Error(err))
		}

		if err = bigQueryClient.InsertTransactions(block.Transactions, ctx); err != nil {
			logger.Error("unable to insert block transactions into BigQuery", zap.Int64("block_number", currentBlock), zap.Error(err))
		}
	}

	return nil
}
