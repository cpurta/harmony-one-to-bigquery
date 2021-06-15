package harmony

//go:generate mockgen -source ./harmony.go -destination ./mock_harmony/mock.go

import "github.com/cpurta/harmony-one-to-bigquery/internal/model"

// HarmonyClient provides a basic interface to access the latest block header submitted
// to the Harmony One blockchain or get all information associated to a specified
// block.
type HarmonyClient interface {
	GetLatestHeader() (*model.Header, error)
	GetBlockByNumber(blockNumber int64) (*model.Block, error)
}
