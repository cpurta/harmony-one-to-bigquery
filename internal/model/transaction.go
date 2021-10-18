package model

import "cloud.google.com/go/bigquery"

func NewRetryTransaction(txn *Transaction, err error) *RetryTransaction {
	return &RetryTransaction{
		Transaction: txn,
		RetryCount:  0,
		Error:       err,
	}
}

type RetryTransaction struct {
	Transaction *Transaction
	RetryCount  int
	Error       error
}

type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	EthHash          string `json:"ethHash"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	R                string `json:"r"`
	S                string `json:"s"`
	ShardID          int64  `json:"shardID"`
	Timestamp        string `json:"timestamp"`
	To               string `json:"to"`
	ToShardID        int64  `json:"toShardID"`
	TransactionIndex string `json:"transactionIndex"`
	V                string `json:"v"`
	Value            string `json:"value"`
}

func (txn *Transaction) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"blockHash":        txn.BlockHash,
		"blockNumber":      txn.BlockNumber,
		"ethHash":          txn.EthHash,
		"from":             txn.From,
		"gas":              txn.Gas,
		"gasPrice":         txn.GasPrice,
		"input":            txn.Input,
		"nonce":            txn.Nonce,
		"r":                txn.R,
		"s":                txn.S,
		"shardID":          txn.ShardID,
		"timestamp":        txn.Timestamp,
		"to":               txn.To,
		"toShardID":        txn.ToShardID,
		"transactionIndex": txn.TransactionIndex,
		"txnHash":          txn.Hash,
		"v":                txn.V,
		"value":            txn.Value,
	}, bigquery.NoDedupeID, nil
}
