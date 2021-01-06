package tron

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/golibs/mock"
	"github.com/trustwallet/golibs/tokentype"
)

var (
	transferSrc, _                 = mock.JsonFromFilePathToString("mocks/" + "transfer.json")
	tokenTransferSrc, _            = mock.JsonFromFilePathToString("mocks/" + "token_transfer.json")
	wantedTransactionsWithToken, _ = mock.JsonFromFilePathToString("mocks/" + "token_txs_response.json")
	wantedTransactionsOnly, _      = mock.JsonFromFilePathToString("mocks/" + "txs_response.json")

	transferDst = blockatlas.Tx{
		ID:     "24a10f7a503e78adc0d7e380b68005531b09e16b9e3f7b524e33f40985d287df",
		Coin:   coin.TRX,
		From:   "TMuA6YqfCeX8EhbfYEg5y7S4DqzSJireY9",
		To:     "TAUN6FwrnwwmaEqYcckffC7wYmbaS6cBiX",
		Fee:    "0", // TODO
		Date:   1564797900,
		Block:  0, // TODO
		Status: blockatlas.StatusCompleted,
		Meta: blockatlas.Transfer{
			Value:    "100666888000000",
			Symbol:   "TRX",
			Decimals: 6,
		},
	}

	tokenTransferDst = blockatlas.Tx{
		ID:     "24a10f7a503e78adc0d7e380b68005531b09e16b9e3f7b524e33f40985d287df",
		Coin:   coin.TRX,
		From:   "TMuA6YqfCeX8EhbfYEg5y7S4DqzSJireY9",
		To:     "TAUN6FwrnwwmaEqYcckffC7wYmbaS6cBiX",
		Fee:    "0", // TODO
		Date:   1564797900,
		Block:  0, // TODO
		Status: blockatlas.StatusCompleted,
		Meta: blockatlas.TokenTransfer{
			Name:     "BitTorrent",
			Symbol:   "BTT",
			TokenID:  "1002000",
			Decimals: 6,
			Value:    "2776267",
			From:     "TMuA6YqfCeX8EhbfYEg5y7S4DqzSJireY9",
			To:       "TAUN6FwrnwwmaEqYcckffC7wYmbaS6cBiX",
		},
	}

	assetInfo = AssetInfo{Name: "BitTorrent", Symbol: "BTT", Decimals: 6, ID: 1002000}
)

type test struct {
	name        string
	apiResponse string
	expected    *blockatlas.Tx
}

func TestNormalizeTokenTransfer(t *testing.T) {
	testNormalizeTokenTransfer(t, &test{
		name:        "token transfer",
		apiResponse: tokenTransferSrc,
		expected:    &tokenTransferDst,
	})
}

func testNormalizeTokenTransfer(t *testing.T, _test *test) {
	var srcTx Tx
	err := json.Unmarshal([]byte(_test.apiResponse), &srcTx)
	assert.NoError(t, err)
	assert.NotNil(t, srcTx)
	res, err := normalize(srcTx)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	addTokenMeta(res, srcTx, assetInfo)
	assert.Equal(t, _test.expected, res)
}

func TestNormalize(t *testing.T) {
	testNormalize(t, &test{
		name:        "transfer",
		apiResponse: transferSrc,
		expected:    &transferDst,
	})
}

func testNormalize(t *testing.T, _test *test) {
	var srcTx Tx
	err := json.Unmarshal([]byte(_test.apiResponse), &srcTx)
	assert.NoError(t, err)
	assert.NotNil(t, srcTx)
	res, err := normalize(srcTx)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, _test.expected, res)
}

func TestPlatform_GetTxsByAddress(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()

	p := Init(server.URL, server.URL)
	res, err := p.GetTxsByAddress("TM1zzNDZD2DPASbKcgdVoTYhfmYgtfwx9R")
	assert.Nil(t, err)

	rawRes, err := json.Marshal(res)
	assert.Nil(t, err)
	assert.JSONEq(t, wantedTransactionsOnly, string(rawRes))
}

func TestPlatform_GetTokenTxsByAddress(t *testing.T) {
	server := httptest.NewServer(createMockedAPI())
	defer server.Close()

	p := Init(server.URL, server.URL)
	res, err := p.GetTokenTxsByAddress("TM1zzNDZD2DPASbKcgdVoTYhfmYgtfwx9D", "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t")
	assert.Nil(t, err)

	rawRes, err := json.Marshal(res)
	assert.Nil(t, err)
	assert.JSONEq(t, wantedTransactionsWithToken, string(rawRes))
}

func Test_getTokenType(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  tokentype.Type
	}{
		{"default trc20", "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", tokentype.TRC20},
		{"default trc10", "1002001", tokentype.TRC10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getTokenType(tt.token))
		})
	}
}
