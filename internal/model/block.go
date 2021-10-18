package model

import "cloud.google.com/go/bigquery"

func NewRetryBlock(block *Block, err error) *RetryBlock {
	return &RetryBlock{
		Block:      block,
		RetryCount: 0,
		Error:      err,
	}
}

type RetryBlock struct {
	Block      *Block
	RetryCount int
	Error      error
}

type Block struct {
	Difficulty       int64          `json:"difficulty"`
	Epoch            string         `json:"epoch"`
	ExtraData        string         `json:"extraData"`
	GasLimit         string         `json:"gasLimit"`
	GasUsed          string         `json:"gasUsed"`
	Hash             string         `json:"hash"`
	LogBloom         string         `json:"logBloom"`
	Miner            string         `json:"miner"`
	MixHash          string         `json:"mixHash"`
	Nonce            int64          `json:"nonce"`
	Number           string         `json:"number"`
	ParentHash       string         `json:"parentHash"`
	ReceiptsRoot     string         `json:"receiptsRoot"`
	Size             string         `json:"size"`
	StateRoot        string         `json:"stateRoot"`
	Timestamp        string         `json:"timestamp"`
	Transactions     []*Transaction `json:"transactions"`
	TransactionsRoot string         `json:"transactionsRoot"`
	ViewID           string         `json:"viewID"`
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
