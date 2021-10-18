package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockUnmarshal(t *testing.T) {
	blockBody := `{
		"jsonrpc": "2.0",
		"id": 1,
		"result": {
		  "number": "0xc65d40",
		  "viewID": "0xc65ea0",
		  "epoch": "0x239",
		  "hash": "0x9bd47aaeb31d7b7b5ed92c69074f90c154a24d9f381b0d5cc1ad3d6ee848970b",
		  "parentHash": "0x72df0ad1be2bfa54f76aa3dac9e37a5e4bfd17c7a44d3bb25c363fbda477e99c",
		  "nonce": 0,
		  "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
		  "logsBloom": "0x00000000000000000000000000001000000000002000000000000000000000000000000000000000000028000000000000000000000010000008000000040000000000000000100000000008000000000000000000040000000000000000000000000000020000000000000000000800000000000000000000000010000000000000100000000000000000000000000000000000000000000000000000000810000000000000000004000000000000000800000000000000000000000000200000040002000000000000000000000000000004000000000000000000000060000004140000000000000000000000000008100000040000000000000000000000",
		  "stateRoot": "0xec910a3b206b632fceb1ee1a410f928eba5438e2c339f27946f09bae4f4728f9",
		  "miner": "one1gh043zc95e6mtutwy5a2zhvsxv7lnlklkj42ux",
		  "difficulty": 0,
		  "extraData": "0x",
		  "size": "0x558",
		  "gasLimit": "0x4c4b400",
		  "gasUsed": "0x21942",
		  "vrf": "0x0000000000000000000000000000000000000000000000000000000000000000",
		  "vrfProof": "0x",
		  "timestamp": "0x609e695c",
		  "transactionsRoot": "0xcd0404df47e2b21f5513264bbf66fdf6d1bf240d0fcf8cb795a96963243c6946",
		  "receiptsRoot": "0x6ecb399e764af286b76c128e271ce35ac0e1c005fd53e7dd6508d72fec23b83b",
		  "uncles": [],
		  "transactions": [
			{
			  "blockHash": "0x9bd47aaeb31d7b7b5ed92c69074f90c154a24d9f381b0d5cc1ad3d6ee848970b",
			  "blockNumber": "0xc65d40",
			  "from": "one1wmudztmxynm38vkc3998fxkeymmczg6st7sf83",
			  "timestamp": "0x609e695c",
			  "gas": "0x5208",
			  "gasPrice": "0x2540be400",
			  "hash": "0x7369a144eaab6f787f306648af14155c92b1bb66be15bf379e94ba74ee1521e1",
			  "ethHash": "0x2d6e9595a7c6a663c283a447a4d40c03f094273601ebc4fe05440c46aa4f6b4d",
			  "input": "0x",
			  "nonce": "0x234ea",
			  "to": "one13hrrej58t6kn3k24fwuhzudy7x9tayh8p73cq9",
			  "transactionIndex": "0x0",
			  "value": "0x107476f644f95e1c",
			  "shardID": 0,
			  "toShardID": 0,
			  "v": "0x25",
			  "r": "0x267b9ff5f6242178903b3e8b746bf0d0607d29b4b40cd1e6601279e2fb7787f3",
			  "s": "0x38fa728f7980c81e4b9f8b2a94bad1c7446aedafb7a9c3f8f17d0fbe4c2fbfe8"
			},
			{
			  "blockHash": "0x9bd47aaeb31d7b7b5ed92c69074f90c154a24d9f381b0d5cc1ad3d6ee848970b",
			  "blockNumber": "0xc65d40",
			  "from": "one1vvlyvzlkv0m3c698e8esljlcavzylmt37ec0hw",
			  "timestamp": "0x609e695c",
			  "gas": "0x30d40",
			  "gasPrice": "0x3b9aca00",
			  "hash": "0xfc0cf5df36e28b9d98992bcc8129ebfe973a42160b91e63c82b3046355d47079",
			  "ethHash": "0x1be313680e2070b4e4523a2da9bb491019d9f70c4cf2d66290e006438167d433",
			  "input": "0x441a3e7000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000",
			  "nonce": "0x1a2",
			  "to": "one143cmv9a93v7vzdk376s3sff0xv06k38u5f4a6n",
			  "transactionIndex": "0x1",
			  "value": "0x0",
			  "shardID": 0,
			  "toShardID": 0,
			  "v": "0xc6ac98a3",
			  "r": "0x35b65bfb28ba654ae61e9858ae0fdb73c2424609e92294f9256e8578ebce169f",
			  "s": "0x3b6187a3a108fc50d19788644384a67c538c91549409b3860e911309925a9ce5"
			}
		  ],
		  "stakingTransactions": []
		}
	  }`

	type blockResponse struct {
		JSONRPC string `json:"jsonrpc"`
		ID      int64  `json:"id"`
		Block   *Block `json:"result"`
	}

	var response blockResponse
	err := json.Unmarshal([]byte(blockBody), &response)
	if err != nil {
		assert.NoError(t, err, "was not expecting an error but recieved one")
	}

	assert.Equal(t, 2, len(response.Block.Transactions))
}
