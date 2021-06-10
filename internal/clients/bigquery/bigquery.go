package bigquery

//go:generate mockgen -source ./bigquery.go -destination ./mock_bigquery/mock.go

import (
	"context"
	"io"

	"github.com/cpurta/harmony-one-to-bigquery/internal/model"
)

type BigQueryClient interface {
	GetMostRecentBlockNumber(ctx context.Context) (int64, error)
	InsertBlock(block *model.Block, ctx context.Context) error
	InsertTransactions(transactions []*model.Transaction, ctx context.Context) error
	io.Closer
}
