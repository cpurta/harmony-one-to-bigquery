package model

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
