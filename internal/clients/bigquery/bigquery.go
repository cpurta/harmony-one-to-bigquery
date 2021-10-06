package bigquery

//go:generate mockgen -source ./bigquery.go -destination ./mock_bigquery/mock.go

import (
	"context"
	"io"

	"github.com/cpurta/harmony-one-to-bigquery/internal/model"
)

type BigQueryClient interface {
	GetMostRecentBlockNumber(ctx context.Context) (int64, error)
	ProjectDatasetExists(ctx context.Context) bool
	CreateBlocksTable(ctx context.Context) error
	BlocksTableExists(ctx context.Context) bool
	InsertBlock(block *model.Block, ctx context.Context) error
	CreateProjectDataset(ctx context.Context) error
	CreateTransactionsTable(ctx context.Context) error
	TransactionsTableExists(ctx context.Context) bool
	InsertTransactions(transactions []*model.Transaction, ctx context.Context) error
	UpdateBlocksSchema(ctx context.Context) error
	UpdateTransactionsSchema(ctx context.Context) error
	io.Closer
}
