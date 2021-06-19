package client

import (
	"errors"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var valid200BlockHeaderResponse = `{
  "jsonrpc": "2.0",
  "id": 0,
  "result": {
    "blockHash": "0xd8d975d2519a9ebf2366f163704df7f3843a017ca9d86707b8222f35067d5039",
    "blockNumber": 14408510,
    "crossLinks": [
      {
        "block-number": 14667721,
        "epoch-number": 612,
        "hash": "0x2235b7b4d06691bb3b62d405ada2d9fa4fe42190e85a449d4e21fc8938595ffe",
        "shard-id": 1,
        "signature": "13ccb54717983eda71e7422ff97f861cd9f8f76c7724ec3c505345e661bb1bef571b63cb8576a62d9878fff2ed704406cf393e722d963b071d43bf3b7f2e02ca21a0eb4fbef47627854434a20e0b2bfc3e63504c8abbddf2e8484dbe9efd6018",
        "signature-bitmap": "ffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffff0f",
        "view-id": 14668006
      },
      {
        "block-number": 14725417,
        "epoch-number": 612,
        "hash": "0xb00b3e72a1396ddcfd1ffa244f33ce61c539deca6a79c537a7789921809ee21a",
        "shard-id": 2,
        "signature": "0740970def3adc6257eb11d774e7bdea3a4bae70f874b8bf7f6fa6d7cc1fd73e9300cf6224996d17e80598e803b3ef06abdf96fe55e307cf079e417cc0a9b3c49c1dc03e4afbe91ccabb6b7bc72a6bb08e43b44fd76d5e1af2aa203606eb2e82",
        "signature-bitmap": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff3f",
        "view-id": 14725664
      },
      {
        "block-number": 14725418,
        "epoch-number": 612,
        "hash": "0xf7f864e143c5e5aad33442317f9eafd649919d88aba6f037887a137c8de6f5a2",
        "shard-id": 2,
        "signature": "07116991bc4bf491a4f5b38ecb6eae5b7384a552b3190be9a99eee7572988a34bd085b7ba3e8dad1a838a73db010ad156f3696f3d554aba2558dde20e9c954b09c6580e246d1ac8251a9b97e4e3cbbb2e790b7920d8bc71b18714cc0194acb00",
        "signature-bitmap": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff3f",
        "view-id": 14725665
      },
      {
        "block-number": 14651759,
        "epoch-number": 612,
        "hash": "0xb268b2e80c2238a51ef17e26d762f827133179eca2f86783f6b39a7ed2f1cfaf",
        "shard-id": 3,
        "signature": "dfd1e2e28369dbcb5efe67cf7f3f4969c6136194a526418076136e20d9914077db1debd69518f91c79021849f6e62f08a68abf89d8b267ff40227941078fc2502d373ab7d8fbb27497eb5eb979ea2da6a7b1fa6e406f68816ffc46afabbe7294",
        "signature-bitmap": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff07",
        "view-id": 14652752
      },
      {
        "block-number": 14651760,
        "epoch-number": 612,
        "hash": "0x038da7f0a9fae64bb3c08639b75ca7a962660a52e14d11fbe86db1caa4a1b833",
        "shard-id": 3,
        "signature": "6e22430debca5e8dd10c12a362b9a2c2f55e58b7dc3840dfd925ffb81e98b043b8d01e9dfd0f0b6e3d6215f69331480cc292691b193e0a5cadc59cceb0f46a0b42075a63375132b61d276db94086dce050774defd349b415c4b8c5210e6a0803",
        "signature-bitmap": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff07",
        "view-id": 14652753
      }
    ],
    "epoch": 612,
    "lastCommitBitmap": "ffffffffffffffffffffffffffffffe3ffffffffffffffffffffffffff07",
    "lastCommitSig": "fd8a2640a3459da2cc46943c8ada55001ba5730288b40c8c418a2f8068181d651954576dff2b0b6ad263495a5a7ee8152283e27924fd85aa6cb935860ae6d3d89f30bd3627cd4cfa79413d6ba76122414517998b5f9a7425f0eca68bbf1fbc07",
    "leader": "one1gh043zc95e6mtutwy5a2zhvsxv7lnlklkj42ux",
    "shardID": 0,
    "timestamp": "2021-06-19 22:01:24 +0000 UTC",
    "unixtime": 1624140084,
    "viewID": 14409079
  }
}`

var valid200BlockResponse = `{
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

var error200BlockResponse = `{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32000,
    "message": "requested block number greater than current block number"
  }
}`

var valid200EpochResponse = `{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x264"
}`

var valid200EpochLastBlockResponse = `{
  "jsonrpc": "2.0",
  "id": 1,
  "result": 14417919
}`

var timeout504Response = `<html><body><p>504 Gateway Timeout</p><body></html>`

func TestGetBlockHeader_200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.s0.t.hmny.io", httpmock.NewStringResponder(200, valid200BlockHeaderResponse))

	harmonyClient := NewHarmonyOneClient("https://api.s0.t.hmny.io", http.DefaultClient, zap.NewNop())

	header, err := harmonyClient.GetLatestHeader()

	assert.NoError(t, err)

	assert.Equal(t, int64(14408510), header.BlockNumber)
}

func TestGetBlockByNumber_200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.s0.t.hmny.io", httpmock.NewStringResponder(200, valid200BlockResponse))

	harmonyClient := NewHarmonyOneClient("https://api.s0.t.hmny.io", http.DefaultClient, zap.NewNop())

	block, err := harmonyClient.GetBlockByNumber(int64(1695732))

	assert.NoError(t, err)

	assert.Equal(t, "0x19dff4", block.Number)
}

func TestGetBlockByNumber_200WithError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.s0.t.hmny.io", httpmock.NewStringResponder(200, error200BlockResponse))

	harmonyClient := NewHarmonyOneClient("https://api.s0.t.hmny.io", http.DefaultClient, zap.NewNop())

	block, err := harmonyClient.GetBlockByNumber(int64(1695732))

	assert.Nil(t, block)

	expectedError := errors.New("received error in header response: [-32000] requested block number greater than current block number")

	assert.Equal(t, expectedError.Error(), err.Error())
}
