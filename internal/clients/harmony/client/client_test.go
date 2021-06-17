package client

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var valid200Response = `{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "difficulty": 0,
    "epoch": "0x53",
    "extraData": "0x",
    "gasLimit": "0x6f05b59d3b20000",
    "gasUsed": "0x0",
    "hash": "0x2b51d8c155a64720a3d738a0a2f1f230129ab3902b803795bd3118a9f57efd86",
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "miner": "one1gh043zc95e6mtutwy5a2zhvsxv7lnlklkj42ux",
    "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "nonce": 0,
    "number": "0x19dff4",
    "parentHash": "0xc2b11434fbb5020d059a8ddd6d9317e38c0b7ca6939076617fc21134985ef289",
    "receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "size": "0x2e5",
    "stakingTransactions": [],
    "stateRoot": "0x842d05eb429b693826a4844ee266ff3bc9ae56a590f281c4cdf0ae8199f7a351",
    "timestamp": "0x5ded21bc",
    "transactions": [],
    "transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "uncles": [],
    "viewID": "0x19e022"
  }
}`

var timeout504Response = `<html><body><p>504 Gateway Timeout</p><body></html>`

func TestGetBlockByNumber_200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.s0.t.hmny.io", httpmock.NewStringResponder(200, valid200Response))

	harmonyClient := NewHarmonyOneClient("https://api.s0.t.hmny.io", http.DefaultClient, zap.NewNop())

	block, err := harmonyClient.GetBlockByNumber(int64(1695732))

	assert.NoError(t, err)

	assert.Equal(t, "0x19dff4", block.Number)
}
