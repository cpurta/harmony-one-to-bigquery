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
	"github.com/cpurta/harmony-one-to-bigquery/internal/schema"
	"github.com/cpurta/harmony-one-to-bigquery/internal/util"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

const MAX_CHAN_LENGTH = 1000

type BackfillRunner struct {
	NodeURL              string
	GoogleCloudProjectID string
	DryRun               bool
	DatasetID            string
	BlocksTableID        string
	TxnsTableID          string
	LoggingProduction    bool
	LogLevel             string
	Concurrency          int
	MaxRetries           int
	StartBlock           int64
	WaitTime             time.Duration
	harmonyClient        harmony.HarmonyClient
	bigQueryClient       bigquery.BigQueryClient
	retryBlockChan       chan *model.RetryBlock
	retryTxnChan         chan *model.RetryTransaction
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

	if runner.LoggingProduction {
		if runner.logger, err = zap.NewProduction(); err != nil {
			return err
		}
	} else {
		if runner.logger, err = zap.NewDevelopment(); err != nil {
			return err
		}
	}

	runner.retryBlockChan = make(chan *model.RetryBlock, MAX_CHAN_LENGTH)
	runner.retryTxnChan = make(chan *model.RetryTransaction, MAX_CHAN_LENGTH)

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
		header  *model.Header
		counter = &counter{
			Lock: &sync.Mutex{},
		}
		blocksTable = fmt.Sprintf("`%s.%s.%s`", runner.GoogleCloudProjectID, runner.DatasetID, runner.BlocksTableID)
		txnsTable   = fmt.Sprintf("`%s.%s.%s`", runner.GoogleCloudProjectID, runner.DatasetID, runner.TxnsTableID)
		wg          = &sync.WaitGroup{}
		err         error
	)

	// dataset check

	runner.logger.Debug("checking if dataset exists")
	if datasetExists := runner.bigQueryClient.ProjectDatasetExists(ctx); !datasetExists {
		runner.logger.Info(fmt.Sprintf("dataset_id \"%s\" does not exists in project_id \"%s\", attempting to create", runner.DatasetID, runner.GoogleCloudProjectID))

		if err = runner.bigQueryClient.CreateProjectDataset(ctx); err != nil {
			runner.logger.Error(fmt.Sprintf("unable to create dataset_id \"%s\" in project_id \"%s\"", runner.DatasetID, runner.GoogleCloudProjectID), zap.Error(err))
			return err
		}
	}

	// transactions table check

	runner.logger.Debug("checking if transactions table exists")
	if txnsExists := runner.bigQueryClient.TransactionsTableExists(ctx); !txnsExists {
		runner.logger.Info(fmt.Sprintf("%s does not exist, attempting to create", txnsTable))

		if err = runner.bigQueryClient.CreateTransactionsTable(ctx); err != nil {
			runner.logger.Error(fmt.Sprintf("unable to create %s table", txnsTable), zap.Error(err))
			return err
		}
	}

	txnsSchema, err := runner.bigQueryClient.GetTransactionSchema(ctx)
	if err != nil {
		return err
	}

	if !util.SchemasEqual(schema.TransactionsTableSchema, txnsSchema) {
		return fmt.Errorf("%s table schema does not match schema.TransactionsTableSchema, please fix schema in GCP console and re-run", txnsTable)
	}

	// blocks table check

	runner.logger.Debug("checking if blocks table exists")
	if blocksExist := runner.bigQueryClient.BlocksTableExists(ctx); !blocksExist {
		runner.logger.Info(fmt.Sprintf("%s does not exist, attempting to create", blocksTable))

		if err = runner.bigQueryClient.CreateBlocksTable(ctx); err != nil {
			runner.logger.Error(fmt.Sprintf("unable to create %s table", blocksTable), zap.Error(err))
			return err
		}
	}

	blocksSchema, err := runner.bigQueryClient.GetBlocksSchema(ctx)
	if err != nil {
		return err
	}

	if !util.SchemasEqual(schema.BlocksTableSchema, blocksSchema) {
		return fmt.Errorf("%s table schema does not match schema.TransactionsTableSchema, please fix schema in GCP console and re-run", blocksTable)
	}

	// blocks and transaction processing

	if header, err = runner.harmonyClient.GetLatestHeader(); err != nil {
		runner.logger.Error("unable to get the most recent block header from the harmony blockchain client", zap.Error(err))
		return err
	}

	runner.logger.Info("received current block header on hmy blockchain", zap.Int64("hmy_block_number", header.BlockNumber))

	if runner.StartBlock != 0 {
		if runner.StartBlock > header.BlockNumber {
			return errors.New("start block provided is greater than the most recent block number in Harmony One blockchain")
		}
		if runner.StartBlock < 0 {
			return errors.New("start block must be a positive number")
		}
	}

	if runner.StartBlock != 0 {
		counter.Count = runner.StartBlock
	} else if counter.Count, err = runner.bigQueryClient.GetMostRecentBlockNumber(ctx); err != nil {
		runner.logger.Error("unable to get the most recent block number stored in BigQuery", zap.Error(err))
		return err
	}

	runner.logger.Info("received most recent block number stored in BigQuery", zap.Int64("bq_block_number", counter.Count))

	for i := 0; i < runner.Concurrency; i++ {
		wg.Add(1)
		runner.logger.Info("spawning backfillBlocks go routine", zap.Int("routine_number", i))
		go runner.backfillBlocks(ctx, counter, wg)
	}

	go runner.retryFailedBlocks()
	go runner.retryFailedTxns()

	wg.Wait()

	return nil
}

func (runner *BackfillRunner) backfillBlocks(ctx context.Context, counter *counter, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		var (
			bigQueryBlock *model.Block
			bigQueryTxns  []*model.Transaction
			header        *model.Header
			block         *model.Block
			err           error
		)

		counter.Lock.Lock()
		currentBlock := counter.Count
		counter.Count += 1
		counter.Lock.Unlock()

		if header, err = runner.harmonyClient.GetLatestHeader(); err != nil {
			runner.logger.Error("unable to get the most recent block header from the harmony blockchain client", zap.Error(err))
			return
		}

		if currentBlock >= header.BlockNumber {
			runner.logger.Info("backfill is up to date with block header, will attempt to wait and retry", zap.Int64("hmy_block_number", header.BlockNumber))
			time.Sleep(runner.WaitTime)
			continue
		}

		blockNumberLogger := runner.logger.With(zap.Int64("block_number", currentBlock))

		if bigQueryBlock, err = runner.bigQueryClient.GetBlock(ctx, currentBlock); err != nil {
			blockNumberLogger.Error("unable to check if block exists in BigQuery", zap.Error(err))
		}

		if bigQueryTxns, err = runner.bigQueryClient.GetTransactions(ctx, currentBlock); err != nil {
			blockNumberLogger.Error("unable to check if block transactions exists in BigQuery", zap.Error(err))
		}

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

		if bigQueryBlock == nil {
			if err = runner.bigQueryClient.InsertBlock(ctx, block); err != nil {
				blockNumberLogger.Error("unable to insert block into BigQuery", zap.Error(err))
				runner.retryBlockChan <- model.NewRetryBlock(block, err)
			}
		}

		if len(block.Transactions) == 0 {
			blockNumberLogger.Info("no transactions in block")
			continue
		}

		numTxns := len(block.Transactions)

		if len(bigQueryTxns) == 0 {
			blockNumberLogger.Debug("inserting transactions", zap.Int("num_txns", numTxns))
			if err = runner.bigQueryClient.InsertTransactions(ctx, block.Transactions); err != nil {
				blockNumberLogger.Error("unable to insert block transactions into BigQuery", zap.Error(err))
				for _, txn := range block.Transactions {
					runner.retryTxnChan <- model.NewRetryTransaction(txn, err)
				}
			}
		}
	}

	runner.logger.Info("backfill block go routine finished")
}

func (runner *BackfillRunner) retryFailedBlocks() {
	runner.logger.Debug("starting routine to retry failed blocks")

	for {
		ctx := context.Background()

		retryBlock := <-runner.retryBlockChan
		if retryBlock.RetryCount >= runner.MaxRetries {
			runner.logger.Error("block was unable to be inserted after max attempts", zap.Error(retryBlock.Error))
			continue
		}

		time.Sleep(time.Duration(retryBlock.RetryCount) * time.Millisecond * 100)

		runner.logger.Debug("attempting to re-insert a failed block")
		if err := runner.bigQueryClient.InsertBlock(ctx, retryBlock.Block); err != nil {
			runner.logger.Debug("retry block insert failed, putting back into retry channel")
			retryBlock.RetryCount++
			retryBlock.Error = err
			runner.retryBlockChan <- retryBlock
		}

	}
}

func (runner *BackfillRunner) retryFailedTxns() {
	runner.logger.Debug("starting routine to retry failed transactions")

	for {
		ctx := context.Background()

		retryTxn := <-runner.retryTxnChan
		if retryTxn.RetryCount >= runner.MaxRetries {
			runner.logger.Error("transaction was unable to be inserted after max attempts", zap.Error(retryTxn.Error))
			continue
		}

		time.Sleep(time.Duration(retryTxn.RetryCount) * time.Millisecond * 100)

		runner.logger.Debug("attempting to re-insert a failed transaction")
		txns := []*model.Transaction{retryTxn.Transaction}
		if err := runner.bigQueryClient.InsertTransactions(ctx, txns); err != nil {
			runner.logger.Debug("retry transaction insert failed, putting back into retry channel")
			retryTxn.RetryCount++
			retryTxn.Error = err
			runner.retryTxnChan <- retryTxn
		}

	}
}
