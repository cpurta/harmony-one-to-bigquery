package client

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	bq "github.com/cpurta/harmony-one-to-bigquery/internal/clients/bigquery"
	"github.com/cpurta/harmony-one-to-bigquery/internal/model"
	"github.com/cpurta/harmony-one-to-bigquery/internal/schema"
	"github.com/cpurta/harmony-one-to-bigquery/internal/util"
	"google.golang.org/api/iterator"
)

var _ bq.BigQueryClient = &bigQueryClient{}

type bigQueryClient struct {
	client        *bigquery.Client
	projectID     string
	datasetID     string
	blocksTableID string
	txnsTableID   string
}

func NewBigQueryClient(ctx context.Context, projectID string, datasetID string, blocksTableID string, txnsTableID string) (*bigQueryClient, error) {
	var (
		bqClient = &bigQueryClient{
			projectID:     projectID,
			datasetID:     datasetID,
			blocksTableID: blocksTableID,
			txnsTableID:   txnsTableID,
		}
		err error
	)

	bqClient.client, err = bigquery.NewClient(ctx, projectID)

	return bqClient, err
}

func (client *bigQueryClient) GetMostRecentBlockNumber(ctx context.Context) (int64, error) {
	var (
		latestBlockNumber = int64(0)
		query             = fmt.Sprintf("SELECT number FROM `%s.%s.%s` ORDER BY timestamp DESC LIMIT 1", client.projectID, client.datasetID, client.blocksTableID)
		bqQuery           = client.client.Query(query)
		job               *bigquery.Job
		status            *bigquery.JobStatus
		it                *bigquery.RowIterator
		err               error
	)

	if job, err = bqQuery.Run(ctx); err != nil {
		return latestBlockNumber, err
	}

	if status, err = job.Wait(ctx); err != nil {
		return latestBlockNumber, err
	}

	if err = status.Err(); err != nil {
		return latestBlockNumber, err
	}

	if it, err = job.Read(ctx); err != nil {
		return latestBlockNumber, err
	}

	for {
		var row []bigquery.Value

		err = it.Next(&row)
		if err == iterator.Done {
			break
		}

		if err != nil {
			return latestBlockNumber, err
		}

		if len(row) != 1 {
			break
		}

		if row[0] == nil {
			break
		}

		if latestBlockNumber, err = util.HexToInt(row[0].(string)); err != nil {
			return latestBlockNumber, err
		}
	}

	return latestBlockNumber, nil
}

func (client *bigQueryClient) ProjectDatasetExists(ctx context.Context) bool {
	dataset := client.client.Dataset(client.datasetID)

	if _, err := dataset.Metadata(ctx); err != nil {
		return false
	}

	return true
}

func (client *bigQueryClient) CreateProjectDataset(ctx context.Context) error {
	dataset := client.client.Dataset(client.datasetID)

	return dataset.Create(ctx, nil)
}

func (client *bigQueryClient) CreateBlocksTable(ctx context.Context) error {
	table := client.client.Dataset(client.datasetID).Table(client.blocksTableID)

	metadata := &bigquery.TableMetadata{
		Schema:         schema.BlocksTableSchema,
		ExpirationTime: time.Now().AddDate(1, 0, 0),
	}

	return table.Create(ctx, metadata)
}

func (client *bigQueryClient) BlocksTableExists(ctx context.Context) bool {
	var (
		table = client.client.Dataset(client.datasetID).Table(client.blocksTableID)
		err   error
	)

	if _, err = table.Metadata(ctx); err != nil {
		return false
	}

	return true
}

func (client *bigQueryClient) InsertBlock(block *model.Block, ctx context.Context) error {
	inserter := client.client.Dataset(client.datasetID).Table(client.blocksTableID).Inserter()

	if err := inserter.Put(ctx, block); err != nil {
		return err
	}

	return nil
}

func (client *bigQueryClient) CreateTransactionsTable(ctx context.Context) error {
	table := client.client.Dataset(client.datasetID).Table(client.txnsTableID)

	metadata := &bigquery.TableMetadata{
		Schema:         schema.TransactionsTableSchema,
		ExpirationTime: time.Now().AddDate(1, 0, 0),
	}

	return table.Create(ctx, metadata)
}

func (client *bigQueryClient) TransactionsTableExists(ctx context.Context) bool {
	var (
		table = client.client.Dataset(client.datasetID).Table(client.txnsTableID)
		err   error
	)

	if _, err = table.Metadata(ctx); err != nil {
		return false
	}

	return true
}

func (client *bigQueryClient) InsertTransactions(transactions []*model.Transaction, ctx context.Context) error {
	inserter := client.client.Dataset(client.datasetID).Table(client.txnsTableID).Inserter()

	if err := inserter.Put(ctx, transactions); err != nil {
		return err
	}

	return nil
}

func (client *bigQueryClient) UpdateBlocksSchema(ctx context.Context) error {
	table := client.client.Dataset(client.datasetID).Table(client.blocksTableID)

	metadata, err := table.Metadata(ctx)
	if err != nil {
		return err
	}

	update := bigquery.TableMetadataToUpdate{
		Schema: schema.BlocksTableSchema,
	}

	if _, err = table.Update(ctx, update, metadata.ETag); err != nil {
		return err
	}

	return nil
}

func (client *bigQueryClient) UpdateTransactionsSchema(ctx context.Context) error {
	table := client.client.Dataset(client.datasetID).Table(client.txnsTableID)

	metadata, err := table.Metadata(ctx)
	if err != nil {
		return err
	}

	update := bigquery.TableMetadataToUpdate{
		Schema: schema.BlocksTableSchema,
	}

	if _, err = table.Update(ctx, update, metadata.ETag); err != nil {
		return err
	}

	return nil
}

func (client *bigQueryClient) GetBlocksSchema(ctx context.Context) (*bigquery.Schema, error) {
	table := client.client.Dataset(client.datasetID).Table(client.blocksTableID)

	metadata, err := table.Metadata(ctx)
	if err != nil {
		return nil, err
	}

	return &metadata.Schema, nil
}

func (client *bigQueryClient) GetTransactionSchema(ctx context.Context) (*bigquery.Schema, error) {
	table := client.client.Dataset(client.datasetID).Table(client.txnsTableID)

	metadata, err := table.Metadata(ctx)
	if err != nil {
		return nil, err
	}

	return &metadata.Schema, nil
}

func (client *bigQueryClient) Close() error {
	return client.client.Close()
}

func hexToInt(hexString string) (int64, error) {
	str := strings.Replace(hexString, "0x", "", -1)
	str = strings.Replace(str, "0X", "", -1)

	return strconv.ParseInt(str, 16, 64)
}
