package schema

import (
	"cloud.google.com/go/bigquery"
)

var (
	BlocksTableSchema = bigquery.Schema{
		{Name: "blockHash", Type: bigquery.StringFieldType},
		{Name: "difficulty", Type: bigquery.StringFieldType},
		{Name: "epoch", Type: bigquery.StringFieldType},
		{Name: "extraData", Type: bigquery.StringFieldType},
		{Name: "gasLimit", Type: bigquery.StringFieldType},
		{Name: "gasUsed", Type: bigquery.StringFieldType},
		{Name: "logBloom", Type: bigquery.StringFieldType},
		{Name: "miner", Type: bigquery.StringFieldType},
		{Name: "mixHash", Type: bigquery.StringFieldType},
		{Name: "nonce", Type: bigquery.IntegerFieldType},
		{Name: "number", Type: bigquery.StringFieldType},
		{Name: "parentHash", Type: bigquery.StringFieldType},
		{Name: "receiptsRoot", Type: bigquery.StringFieldType},
		{Name: "size", Type: bigquery.StringFieldType},
		{Name: "stateRoot", Type: bigquery.StringFieldType},
		{Name: "timestamp", Type: bigquery.StringFieldType},
		{Name: "transactionsRoot", Type: bigquery.StringFieldType},
		{Name: "viewID", Type: bigquery.StringFieldType},
	}
	TransactionsTableSchema = bigquery.Schema{
		{Name: "blockHash", Type: bigquery.StringFieldType},
		{Name: "blockNumber", Type: bigquery.StringFieldType},
		{Name: "ethHash", Type: bigquery.StringFieldType},
		{Name: "from", Type: bigquery.StringFieldType},
		{Name: "gas", Type: bigquery.StringFieldType},
		{Name: "gasPrice", Type: bigquery.StringFieldType},
		{Name: "input", Type: bigquery.StringFieldType},
		{Name: "nonce", Type: bigquery.StringFieldType},
		{Name: "r", Type: bigquery.StringFieldType},
		{Name: "s", Type: bigquery.StringFieldType},
		{Name: "shardID", Type: bigquery.IntegerFieldType},
		{Name: "timestamp", Type: bigquery.StringFieldType},
		{Name: "to", Type: bigquery.StringFieldType},
		{Name: "toShardID", Type: bigquery.IntegerFieldType},
		{Name: "transactionIndex", Type: bigquery.StringFieldType},
		{Name: "txnHash", Type: bigquery.StringFieldType},
		{Name: "v", Type: bigquery.StringFieldType},
		{Name: "value", Type: bigquery.StringFieldType},
	}
)
