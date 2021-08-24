package model

import "cloud.google.com/go/bigquery"

type Block struct {
	Difficulty       int64
	Epoch            string
	ExtraData        string
	GasLimit         string
	GasUsed          string
	Hash             string
	LogBloom         string
	Miner            string
	MixHash          string
	Nonce            int64
	Number           string
	ParentHash       string
	ReceiptsRoot     string
	Size             string
	StateRoot        string
	Timestamp        string
	Transactions     []*Transaction
	TransactionsRoot string
	ViewID           string
}

func (block *Block) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"blockHash":        block.Hash,
		"difficulty":       block.Difficulty,
		"epoch":            block.Epoch,
		"extraData":        block.ExtraData,
		"gasLimit":         block.GasLimit,
		"gasUsed":          block.GasUsed,
		"logBloom":         block.LogBloom,
		"miner":            block.Miner,
		"mixHash":          block.MixHash,
		"nonce":            block.Nonce,
		"number":           block.Number,
		"parentHash":       block.ParentHash,
		"receiptsRoot":     block.ReceiptsRoot,
		"size":             block.Size,
		"stateRoot":        block.StateRoot,
		"timestamp":        block.Timestamp,
		"transactionsRoot": block.TransactionsRoot,
		"viewID":           block.ViewID,
	}, bigquery.NoDedupeID, nil
}
