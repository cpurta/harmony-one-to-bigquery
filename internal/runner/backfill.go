package runner

import (
	"github.com/cpurta/harmony-one-to-bigquery/internal/clients/bigquery"
	"github.com/cpurta/harmony-one-to-bigquery/internal/clients/harmony"
	"github.com/cpurta/harmony-one-to-bigquery/internal/model"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type BackfillRunner struct {
	HarmonyClient  harmony.HarmonyClient
	BigQueryClient bigquery.BigQueryClient
	Logger         *zap.Logger
}

func (runner *BackfillRunner) Run(cli *cli.Context) error {
	var (
		header       *model.Header
		currentBlock int64
		backfillUpTo int64
		err          error
	)

	if header, err = runner.HarmonyClient.GetLatestHeader(); err != nil {
		runner.Logger.Error("unable to get the most recent block header from the harmony blockchain client", zap.Error(err))
		return err
	}

	if backfillUpTo, err = runner.BigQueryClient.GetMostRecentBlockNumber(); err != nil {
		runner.Logger.Error("unable to get the most recent block number stored in BigQuery", zap.Error(err))
		return err
	}

	currentBlock = header.BlockNumber

	for currentBlock > backfillUpTo {
		var block *model.Block

		if block, err = runner.HarmonyClient.GetBlockByNumber(currentBlock); err != nil {
			runner.Logger.Error("unable get block from harmony blockchain client", zap.Int64("block_number", currentBlock), zap.Error(err))
			continue
		}

		if err = runner.BigQueryClient.InsertBlock(block); err != nil {
			runner.Logger.Error("unable to insert block into BigQuery", zap.Int64("block_number", currentBlock), zap.Error(err))
		}

		if err = runner.BigQueryClient.InsertTransactions(block.Transactions); err != nil {
			runner.Logger.Error("unable to insert block transactions into BigQuery", zap.Int64("block_number", currentBlock), zap.Error(err))
		}

		currentBlock -= 1
	}

	return nil
}
