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
	Difficulty       int64          `json:"difficulty" bigquery:"difficulty"`
	Epoch            string         `json:"epoch" bigquery:"epoch"`
	ExtraData        string         `json:"extraData" bigquery:"extraData"`
	GasLimit         string         `json:"gasLimit" bigquery:"gasLimit"`
	GasUsed          string         `json:"gasUsed" bigquery:"gasUsed"`
	Hash             string         `json:"hash" bigquery:"hash"`
	LogBloom         string         `json:"logBloom" bigquery:"logBloom"`
	Miner            string         `json:"miner" bigquery:"miner"`
	MixHash          string         `json:"mixHash" bigquery:"mixHash"`
	Nonce            int64          `json:"nonce" bigquery:"nonce"`
	Number           string         `json:"number" bigquery:"number"`
	ParentHash       string         `json:"parentHash" bigquery:"parentHash"`
	ReceiptsRoot     string         `json:"receiptsRoot" bigquery:"receiptsRoot"`
	Size             string         `json:"size" bigquery:"size"`
	StateRoot        string         `json:"stateRoot" bigquery:"stateRoot"`
	Timestamp        string         `json:"timestamp" bigquery:"timestamp"`
	Transactions     []*Transaction `json:"transactions" bigquery:"transactions"`
	TransactionsRoot string         `json:"transactionsRoot" bigquery:"transactionsRoot"`
	ViewID           string         `json:"viewID" bigquery:"viewID"`
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
