package zilliqa

import (
	"reflect"
	"testing"

	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/golibs/mock"
	"github.com/trustwallet/golibs/txtype"
)

func TestNormalizeTx(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantTx  txtype.Tx
		wantErr bool
	}{
		{
			name: "Test normalize transaction",
			args: args{
				filename: "transfer.json",
			},
			wantTx: txtype.Tx{
				ID:       "0xd44413c79e7518152f3b05ef1edff8ef59afd06119b16d09c8bc72e94fed7843",
				Coin:     coin.ZIL,
				From:     "0x88af5ba10796d9091d6893eed4db23ef0bbbca37",
				To:       "0x7fccacf066a5f26ee3affc2ed1fa9810deaa632c",
				Fee:      "1000000000",
				Date:     1557889788,
				Block:    104282,
				Status:   txtype.StatusCompleted,
				Sequence: 3,
				Memo:     "",
				Meta: txtype.Transfer{
					Value:    "7997000000000",
					Symbol:   "ZIL",
					Decimals: 12,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var srcTx Tx
			_ = mock.JsonModelFromFilePath("mocks/"+tt.args.filename, &srcTx)
			gotTx := Normalize(&srcTx)
			if !reflect.DeepEqual(gotTx, tt.wantTx) {
				t.Errorf("NormalizeTx() gotTx = %v, want %v", gotTx, tt.wantTx)
			}
		})
	}
}
