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
}

func (runner *BackfillRunner) Run(cli *cli.Context) error {
	var (
		header         *model.Header
		harmonyClient  harmony.HarmonyClient
		bigQueryClient bigquery.BigQueryClient
		currentBlock   int64
		backfillUpTo   = int64(0)
		logger         *zap.Logger
		ctx            = context.Background()
		err            error
	)

	if logger, err = zap.NewDevelopment(); err != nil {
		return err
	}

	harmonyClient = client.NewHarmonyOneClient(runner.NodeURL, http.DefaultClient)

	if bigQueryClient, err = bq.NewBigQueryClient(ctx, runner.GoogleCloudProjectID, runner.DatasetID, runner.BlocksTableID, runner.TxnsTableID); err != nil {
		return err
	}

	if header, err = harmonyClient.GetLatestHeader(); err != nil {
		logger.Error("unable to get the most recent block header from the harmony blockchain client", zap.Error(err))
		return err
	}

	if backfillUpTo, err = bigQueryClient.GetMostRecentBlockNumber(ctx); err != nil {
		logger.Error("unable to get the most recent block number stored in BigQuery", zap.Error(err))
		return err
	}

	currentBlock = header.BlockNumber

	for ; currentBlock > backfillUpTo; currentBlock-- {
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
