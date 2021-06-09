package model

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
