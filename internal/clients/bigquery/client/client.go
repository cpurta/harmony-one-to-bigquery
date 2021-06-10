package client

import (
	"context"

	"cloud.google.com/go/bigquery"
	bq "github.com/cpurta/harmony-one-to-bigquery/internal/clients/bigquery"
	"github.com/cpurta/harmony-one-to-bigquery/internal/model"
	"google.golang.org/api/iterator"
)

var _ bq.BigQueryClient = &bigQueryClient{}

type bigQueryClient struct {
	client *bigquery.Client
}

func NewBigQueryClient(ctx context.Context, projectID string) (*bigQueryClient, error) {
	var (
		bqClient = &bigQueryClient{}
		err      error
	)

	bqClient.client, err = bigquery.NewClient(ctx, projectID)

	return bqClient, err
}

func (client *bigQueryClient) GetMostRecentBlockNumber(ctx context.Context) (int64, error) {
	latestBlockNumber = int64(0)

	query := client.client.Query("SELECT MAX(block_number) FROM `bigquery-public-data.` LIMIT 1")

	job, err := query.Run(ctx)
	if err != nil {
		return int64(-1), err
	}

	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}
	it, err := job.Read(ctx)
	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		latestBlockNumber = row.(int64)
	}

	return latestBlockNumber, nil
}

func (client *bigQueryClient) InsertBlock(block *model.Block, ctx context.Context) error {
	return nil
}

func (client *bigQueryClient) InsertTransactions(transactions []*model.Transaction, ctx context.Context) error {
	return nil
}

func (client *bigQueryClient) Close() {
	client.client.Close()
}
