package bigquery

//go:generate mockgen -source ./bigquery.go -destination ./mock_bigquery/mock.go

import "github.com/cpurta/harmony-one-to-bigquery/internal/model"

type BigQueryClient interface {
	GetMostRecentBlockNumber() (int64, error)
	InsertBlock(block *model.Block) error
	InsertTransactions(transactions []*model.Transaction) error
}
