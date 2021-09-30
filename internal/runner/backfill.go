package runner

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
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
	Concurrency          int
	harmonyClient        harmony.HarmonyClient
	bigQueryClient       bigquery.BigQueryClient
	logger               *zap.Logger
}

type counter struct {
	Count int64
	Lock  *sync.Mutex
}

func (runner *BackfillRunner) Run(cli *cli.Context) error {
	var (
		ctx = context.Background()
		err error
	)

	if runner.logger, err = zap.NewProduction(); err != nil {
		return err
	}

	if runner.Concurrency <= 0 {
		return errors.New("must specify concurrency > 0")
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
		backfillUpTo = int64(0)
		counter      = &counter{
			Lock: &sync.Mutex{},
		}
		blocksTable = fmt.Sprintf("`%s.%s.%s`", runner.GoogleCloudProjectID, runner.DatasetID, runner.BlocksTableID)
		txnsTable   = fmt.Sprintf("`%s.%s.%s`", runner.GoogleCloudProjectID, runner.DatasetID, runner.TxnsTableID)
		wg          = &sync.WaitGroup{}
		err         error
	)

	runner.logger.Debug("checking if dataset exists")
	if datasetExists := runner.bigQueryClient.ProjectDatasetExists(ctx); !datasetExists {
		runner.logger.Info(fmt.Sprintf("dataset_id \"%s\" does not exists in project_id \"%s\", attempting to create", runner.DatasetID, runner.GoogleCloudProjectID))

		if err = runner.bigQueryClient.CreateProjectDataset(ctx); err != nil {
			runner.logger.Error(fmt.Sprintf("unable to create dataset_id \"%s\" in project_id \"%s\"", runner.DatasetID, runner.GoogleCloudProjectID), zap.Error(err))
			return err
		}
	}

	runner.logger.Debug("checking if transactions table exists")
	if txnsExists := runner.bigQueryClient.BlocksTableExists(ctx); !txnsExists {
		runner.logger.Info(fmt.Sprintf("%s does not exist, attempting to create", txnsTable))

		if err = runner.bigQueryClient.CreateTransactionsTable(ctx); err != nil {
			runner.logger.Error(fmt.Sprintf("unable to create %s table", txnsTable), zap.Error(err))
			return err
		}

		runner.logger.Debug("cold start: waiting until transactions tables are recognized")
		time.Sleep(time.Second * 2)
	}

	runner.logger.Debug("checking if blocks table exists")
	if blocksExist := runner.bigQueryClient.BlocksTableExists(ctx); !blocksExist {
		runner.logger.Info(fmt.Sprintf("%s does not exist, attempting to create", blocksTable))

		if err = runner.bigQueryClient.CreateBlocksTable(ctx); err != nil {
			runner.logger.Error(fmt.Sprintf("unable to create %s table", blocksTable), zap.Error(err))
			return err
		}

		runner.logger.Debug("cold start: waiting until blocks tables are recognized")
		time.Sleep(time.Second * 2)
	}

	if header, err = runner.harmonyClient.GetLatestHeader(); err != nil {
		runner.logger.Error("unable to get the most recent block header from the harmony blockchain client", zap.Error(err))
		return err
	}

	runner.logger.Info("received current block header on hmy blockchain", zap.Int64("hmy_block_number", header.BlockNumber))

	if counter.Count, err = runner.bigQueryClient.GetMostRecentBlockNumber(ctx); err != nil {
		runner.logger.Error("unable to get the most recent block number stored in BigQuery", zap.Error(err))
		return err
	}

	runner.logger.Info("received most recent block number stored in BigQuery", zap.Int64("bq_block_number", counter.Count))

	backfillUpTo = header.BlockNumber

	for i := 0; i < runner.Concurrency; i++ {
		wg.Add(1)
		runner.logger.Info("spawning backfillBlocks go routine", zap.Int("routine_number", i))
		go runner.backfillBlocks(ctx, counter, wg, backfillUpTo)
	}

	wg.Wait()

	return nil
}

func (runner *BackfillRunner) backfillBlocks(ctx context.Context, counter *counter, wg *sync.WaitGroup, endBlock int64) {
	defer wg.Done()

	for {
		var (
			block *model.Block
			err   error
		)

		counter.Lock.Lock()
		currentBlock := counter.Count
		counter.Count += 1
		counter.Lock.Unlock()

		if currentBlock >= endBlock {
			break
		}

		blockNumberLogger := runner.logger.With(zap.Int64("block_number", currentBlock))

		if block, err = runner.harmonyClient.GetBlockByNumber(currentBlock); err != nil {
			if strings.Contains(err.Error(), "-32000") {
				blockNumberLogger.Info("recieve error stating we are greater than current block, bailing...")
				break
			}

			blockNumberLogger.Error("unable get block from harmony blockchain client", zap.Error(err))
			continue
		}

		if block == nil {
			continue
		}

		if currentBlock%1000 == 0 && currentBlock != 0 {
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

	runner.logger.Info("backfill block go routine finished")
}
