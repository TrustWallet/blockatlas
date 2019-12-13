package cosmos

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"testing"

	"github.com/trustwallet/blockatlas/coin"
)

const transferSrc = `
{
	"height": "151980",
	"txhash": "E19B011D20D862DA0BEA7F24E3BC6DFF666EE6E044FCD9BD95B073478086DBB6",
	"raw_log": "[{\"msg_index\":\"0\",\"success\":true,\"log\":\"\"}]",
	"logs": [
	  {
		"msg_index": "0",
		"success": true,
		"log": ""
	  }
	],
	"gas_wanted": "100000",
	"gas_used": "27678",
	"tags": [
	  {
		"key": "action",
		"value": "send"
	  },
	  {
		"key": "sender",
		"value": "cosmos1rw62phusuv9vzraezr55k0vsqssvz6ed52zyrl"
	  },
	  {
		"key": "recipient",
		"value": "cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae"
	  }
	],
	"tx": {
	  "type": "auth/StdTx",
	  "value": {
		"msg": [
		  {
			"type": "cosmos-sdk/MsgSend",
			"value": {
			  "from_address": "cosmos1rw62phusuv9vzraezr55k0vsqssvz6ed52zyrl",
			  "to_address": "cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae",
			  "amount": [
				{
				  "denom": "uatom",
				  "amount": "2271999999"
				}
			  ]
			}
		  }
		],
		"fee": {
		  "amount": [
			{
			  "denom": "uatom",
			  "amount": "1"
			}
		  ],
		  "gas": "100000"
		},
		"signatures": [
		  {
			"pub_key": {
			  "type": "tendermint/PubKeySecp256k1",
			  "value": "A21fdP6IbVC9hER5smiim8I4EbFeIF/bW81IKwmmsdjH"
			},
			"signature": "MuR85p714L94tCenogRqzLh1bsbmhKTjs1L9JJPdhSVwQKh61EGlLqYGoUeN/n9xb+OOR9ESUOh2CAzVulKoVQ=="
		  }
		],
		"memo": ""
	  }
	},
	"timestamp": "2019-05-04T17:57:57Z"
  }
`

const failedTransferSrc = `
{
  "height": "5552",
  "txhash": "5E78C65A8C1A6C8239EBBBBF2E42020E6ADBA8037EDEA83BF88E1A9159CF13B8",
  "code": 12,
  "raw_log": "{\"codespace\":\"sdk\",\"code\":12,\"message\":\"out of gas in location: WritePerByte; gasWanted: 40000, gasUsed: 40480\"}",
  "gas_wanted": "40000",
  "gas_used": "40480",
  "tx": {
    "type": "cosmos-sdk/StdTx",
    "value": {
      "msg": [
        {
          "type": "cosmos-sdk/MsgSend",
          "value": {
            "from_address": "cosmos1shpfyt7psrff2ux7nznxvj6f7gq59fcqng5mku",
            "to_address": "cosmos1za4pu5gxm80fg6sx0956f88l2sx7jfg2vf7nlc",
            "amount": [
              {
                "denom": "uatom",
                "amount": "100000"
              }
            ]
          }
        }
      ],
      "fee": {
        "amount": [
          {
            "denom": "uatom",
            "amount": "2000"
          }
        ],
        "gas": "40000"
      },
      "signatures": [
        {
          "pub_key": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "A0IDIokqw01U2YcdylvqD/sJHW5w9puS5vZWSf2GUaqL"
          },
          "signature": "1Kwp4dBZUbVV6Fk8AFcmNfSqi7MXFfqyLvHexFZXoqcKh+sNuezry89RhDAWgSMNLyaK20hI2XcUyks+Vo4QEQ=="
        }
      ],
      "memo": "UniCoins registration rewards"
    }
  },
  "timestamp": "2019-12-12T03:21:42Z"
}`

const delegateSrc = `
{  
   "height":"1258202",
   "txhash":"11078091D1D5BD84F4275B6CEE02170428944DB0E8EEC37E980551435F6D04C7",
   "raw_log":"[{\"msg_index\":\"0\",\"success\":true,\"log\":\"\"}]",
   "logs":[  
      {  
         "msg_index":"0",
         "success":true,
         "log":""
      }
   ],
   "gas_wanted":"200000",
   "gas_used":"103206",
   "tags":[  
      {  
         "key":"action",
         "value":"delegate"
      },
      {  
         "key":"delegator",
         "value":"cosmos1237l0vauhw78qtwq045jd24ay4urpec6r3xfn3"
      },
      {  
         "key":"destination-validator",
         "value":"cosmosvaloper12w6tynmjzq4l8zdla3v4x0jt8lt4rcz5gk7zg2"
      }
   ],
   "tx":{  
      "type":"auth/StdTx",
      "value":{  
         "msg":[  
            {  
               "type":"cosmos-sdk/MsgDelegate",
               "value":{  
                  "delegator_address":"cosmos1237l0vauhw78qtwq045jd24ay4urpec6r3xfn3",
                  "validator_address":"cosmosvaloper12w6tynmjzq4l8zdla3v4x0jt8lt4rcz5gk7zg2",
                  "amount":{  
                     "denom":"uatom",
                     "amount":"49920"
                  }
               }
            }
         ],
         "fee":{  
            "amount":[  
               {  
                  "denom":"uatom",
                  "amount":"5000"
               }
            ],
            "gas":"200000"
         },
         "signatures":[  
            {  
               "pub_key":{  
                  "type":"tendermint/PubKeySecp256k1",
                  "value":"AsZL4GaIEGW6ogh1rEasxHtmirpeBnycLz4VR0rSVr9p"
               },
               "signature":"w6sNVzTSsE32ERbBdYYySSp6nj+4xNODuq5GKRVb8q04jMHUbx9AhuZeAhYrkvdkzOl3bD7vRYGx9P1V6yHj0A=="
            }
         ],
         "memo":""
      }
   },
   "timestamp":"2019-08-01T04:10:16Z"
}
`

const unDelegateSrc = `
{  
   "height":"1257037",
   "txhash":"A1EC36741FEF681F4A77B8F6032AD081100EE5ECB4CC76AEAC2174BC6B871CFE",
   "data":"0C0889ECF7EA0510FB9D8CAD03",
   "raw_log":"[{\"msg_index\":\"0\",\"success\":true,\"log\":\"\"}]",
   "logs":[  
      {  
         "msg_index":"0",
         "success":true,
         "log":""
      }
   ],
   "gas_wanted":"200000",
   "gas_used":"107804",
   "tags":[  
      {  
         "key":"action",
         "value":"begin_unbonding"
      },
      {  
         "key":"delegator",
         "value":"cosmos137rrp4p8n0nqcft0mwc62tdnyhhzf80knv5t94"
      },
      {  
         "key":"source-validator",
         "value":"cosmosvaloper1te8nxpc2myjfrhaty0dnzdhs5ahdh5agzuym9v"
      },
      {  
         "key":"end-time",
         "value":"2019-08-22T01:55:21Z"
      }
   ],
   "tx":{  
      "type":"auth/StdTx",
      "value":{  
         "msg":[  
            {  
               "type":"cosmos-sdk/MsgUndelegate",
               "value":{  
                  "delegator_address":"cosmos137rrp4p8n0nqcft0mwc62tdnyhhzf80knv5t94",
                  "validator_address":"cosmosvaloper1te8nxpc2myjfrhaty0dnzdhs5ahdh5agzuym9v",
                  "amount":{  
                     "denom":"uatom",
                     "amount":"5100000000"
                  }
               }
            }
         ],
         "fee":{  
            "amount":[  
               {  
                  "denom":"uatom",
                  "amount":"5000"
               }
            ],
            "gas":"200000"
         },
         "signatures":[  
            {  
               "pub_key":{  
                  "type":"tendermint/PubKeySecp256k1",
                  "value":"A+tPzMXCW7vxmW5VN9Q/CO+fxnEXYlSMOklDVgaFutQD"
               },
               "signature":"rh25A/RTm8TUTUGOGhufqxn9vLFef/04xEKMJLUD5QhBVabRADvEgAP1J842XTDtVBS0SpVD/MrPduqRp0nNzg=="
            }
         ],
         "memo":""
      }
   },
   "timestamp":"2019-08-01T01:55:21Z"
}
`

var transferDst = blockatlas.Tx{
	ID:     "E19B011D20D862DA0BEA7F24E3BC6DFF666EE6E044FCD9BD95B073478086DBB6",
	Coin:   coin.ATOM,
	From:   "cosmos1rw62phusuv9vzraezr55k0vsqssvz6ed52zyrl",
	To:     "cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae",
	Fee:    "1",
	Date:   1556992677,
	Block:  151980,
	Status: blockatlas.StatusCompleted,
	Type:   blockatlas.TxTransfer,
	Meta: blockatlas.Transfer{
		Value:    "2271999999",
		Symbol:   "ATOM",
		Decimals: 6,
	},
}

var delegateDst = blockatlas.Tx{
	ID:     "11078091D1D5BD84F4275B6CEE02170428944DB0E8EEC37E980551435F6D04C7",
	Coin:   coin.ATOM,
	From:   "cosmos1237l0vauhw78qtwq045jd24ay4urpec6r3xfn3",
	To:     "cosmosvaloper12w6tynmjzq4l8zdla3v4x0jt8lt4rcz5gk7zg2",
	Fee:    "5000",
	Date:   1564632616,
	Block:  1258202,
	Status: blockatlas.StatusCompleted,
	Type:   blockatlas.TxAnyAction,
	Meta: blockatlas.AnyAction{
		Coin:     coin.ATOM,
		Title:    blockatlas.AnyActionDelegation,
		Key:      blockatlas.KeyStakeDelegate,
		Name:     "ATOM",
		Symbol:   coin.Coins[coin.ATOM].Symbol,
		Decimals: coin.Coins[coin.ATOM].Decimals,
		Value:    "49920",
	},
}

var unDelegateDst = blockatlas.Tx{
	ID:     "A1EC36741FEF681F4A77B8F6032AD081100EE5ECB4CC76AEAC2174BC6B871CFE",
	Coin:   coin.ATOM,
	From:   "cosmos137rrp4p8n0nqcft0mwc62tdnyhhzf80knv5t94",
	To:     "cosmosvaloper1te8nxpc2myjfrhaty0dnzdhs5ahdh5agzuym9v",
	Fee:    "5000",
	Date:   1564624521,
	Block:  1257037,
	Status: blockatlas.StatusCompleted,
	Type:   blockatlas.TxAnyAction,
	Meta: blockatlas.AnyAction{
		Coin:     coin.ATOM,
		Title:    blockatlas.AnyActionUndelegation,
		Key:      blockatlas.KeyStakeDelegate,
		Name:     "ATOM",
		Symbol:   coin.Coins[coin.ATOM].Symbol,
		Decimals: coin.Coins[coin.ATOM].Decimals,
		Value:    "5100000000",
	},
}
var failedTransferDst = blockatlas.Tx{
	ID:     "5E78C65A8C1A6C8239EBBBBF2E42020E6ADBA8037EDEA83BF88E1A9159CF13B8",
	Coin:   coin.ATOM,
	From:   "cosmos1shpfyt7psrff2ux7nznxvj6f7gq59fcqng5mku",
	To:     "cosmos1za4pu5gxm80fg6sx0956f88l2sx7jfg2vf7nlc",
	Fee:    "2000",
	Date:   1576120902,
	Block:  5552,
	Status: blockatlas.StatusFailed,
	Type:   blockatlas.TxTransfer,
	Memo:   "UniCoins registration rewards",
	Meta: blockatlas.Transfer{
		Value:    "100000",
		Symbol:   "ATOM",
		Decimals: 6,
	},
}

func TestNormalize(t *testing.T) {
	testNormalize(t, "test transfer tx", transferSrc, transferDst)
	testNormalize(t, "test delegate tx", delegateSrc, delegateDst)
	testNormalize(t, "test undelegate tx", unDelegateSrc, unDelegateDst)
	testNormalize(t, "test failed tx", failedTransferSrc, failedTransferDst)
}

func testNormalize(t *testing.T, name, src string, dst blockatlas.Tx) {
	t.Run(name, func(t *testing.T) {
		var srcTx Tx
		err := json.Unmarshal([]byte(src), &srcTx)
		assert.Nil(t, err)
		tx, ok := Normalize(&srcTx)
		assert.True(t, ok)
		assert.Equal(t, dst, tx, "transfer: tx don't equal")
	})
}
