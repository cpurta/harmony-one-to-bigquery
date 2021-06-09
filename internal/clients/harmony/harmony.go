package harmony

//go:generate mockgen -source ./harmony.go -destination ./mock_harmony/mock.go

import "github.com/cpurta/harmony-one-to-bigquery/internal/model"

type HarmonyClient interface {
	GetLatestHeader() (*model.Header, error)
	GetBlockByNumber(blockNumber int64) (*model.Block, error)
}
