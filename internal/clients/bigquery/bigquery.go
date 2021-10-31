package bigquery

//go:generate mockgen -source ./bigquery.go -destination ./mock_bigquery/mock.go

import (
	"context"
	"io"

	"cloud.google.com/go/bigquery"
	"github.com/cpurta/harmony-one-to-bigquery/internal/model"
)

type BigQueryClient interface {
	GetMostRecentBlockNumber(ctx context.Context) (int64, error)
	ProjectDatasetExists(ctx context.Context) bool
	CreateBlocksTable(ctx context.Context) error
	BlocksTableExists(ctx context.Context) bool
	InsertBlock(ctx context.Context, block *model.Block) error
	CreateProjectDataset(ctx context.Context) error
	CreateTransactionsTable(ctx context.Context) error
	TransactionsTableExists(ctx context.Context) bool
	InsertTransactions(ctx context.Context, transactions []*model.Transaction) error
	GetBlocksSchema(ctx context.Context) (bigquery.Schema, error)
	GetTransactionSchema(ctx context.Context) (bigquery.Schema, error)
	GetBlock(ctx context.Context, blockNumber int64) (*model.Block, error)
	GetTransactions(ctx context.Context, blockNumber int64) ([]*model.Transaction, error)
	io.Closer
}
