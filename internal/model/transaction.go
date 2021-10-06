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
	BlockHash        string
	BlockNumber      string
	EthHash          string
	From             string
	Gas              string
	GasPrice         string
	Hash             string
	Input            string
	Nonce            string
	R                string
	S                string
	ShardID          int64
	Timestamp        string
	To               string
	ToShardID        int64
	TransactionIndex string
	V                string
	Value            string
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
