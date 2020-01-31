package zilliqa

import (
	"encoding/json"
	"reflect"
	"testing"
)

const transaction = `{
	"ID": "f73cf0a229a3d71e1a5c2ac4acbab598c706e64882a2e7c5ed6e406ce69fc16c",
	"amount": "1380000000000",
	"gasLimit": "1",
	"gasPrice": "1000000000",
	"nonce": "16109",
	"receipt": {
	  "cumulative_gas": "1",
	  "epoch_num": "185343",
	  "success": true
	},
	"senderPubKey": "0x02025E984E9FD5ED78537765735C011124A49F2F7543683884FAA685ABC2D3ADC4",
	"signature": "0xF165643EA12514F62297854CE14F2C4EEFE0E19670A6A64E3C497E19442D0B36A91A8790FE320EC48DDCD3E212F0863955FB6AF5436422461916319D5133886D",
	"toAddr": "619a0c9716aef2bc84aafd7ee56e5c2af4e62325",
	"version": "65537"
}`

func TestTxRPC_toTx(t *testing.T) {
	tx := Tx{
		Hash:           "0xf73cf0a229a3d71e1a5c2ac4acbab598c706e64882a2e7c5ed6e406ce69fc16c",
		BlockHeight:    185343,
		From:           "zil1jrpjd8pjuv50cfkfr7eu6yrm3rn5u8rulqhqpz",
		To:             "zil1vxdqe9ck4metep92l4lw2mju9t6wvge9zwkyyl",
		Value:          "1380000000000",
		Fee:            "1000000000",
		Signature:      "0xF165643EA12514F62297854CE14F2C4EEFE0E19670A6A64E3C497E19442D0B36A91A8790FE320EC48DDCD3E212F0863955FB6AF5436422461916319D5133886D",
		Nonce:          "16109",
		ReceiptSuccess: true,
	}

	var txRPC TxRPC
	err := json.Unmarshal([]byte(transaction), &txRPC)
	if err != nil {
		t.Error(err)
		return
	}

	if got := txRPC.toTx(); !reflect.DeepEqual(got, tx) {
		t.Errorf("TxRPC.toTx() = %v, want %v", got, tx)
	}
}

func TestTx_NonceValue(t *testing.T) {
	tests := []struct {
		name  string
		nonce interface{}
		want  uint64
	}{
		{"test int", 0, 0},
		{"test float", 3.4, 3},
		{"test string", "33", 33},
		{"test error string", "test", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := Tx{
				Nonce: tt.nonce,
			}
			if got := tx.NonceValue(); got != tt.want {
				t.Errorf("NonceValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
