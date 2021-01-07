package elrond

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/golibs/mock"
)

var (
	userAddress       = `erd10yagg2vme2jns9zqf9xn8kl86fkc6dr063vnuj0mz2kk2jw0qwuqmfmaw0`
	txTransferSrc1, _ = mock.JsonFromFilePathToString("mocks/tx.json")
	txTransferSrc2, _ = mock.JsonFromFilePathToString("mocks/tx_2.json")
	txTransferSrc3, _ = mock.JsonFromFilePathToString("mocks/tx_3.json")
	txTransferSrc4, _ = mock.JsonFromFilePathToString("mocks/tx_4.json")
	txTransferSrc5, _ = mock.JsonFromFilePathToString("mocks/tx_5.json")
	txTransferSrc6, _ = mock.JsonFromFilePathToString("mocks/tx_6.json")

	txTransfer1Normalized = blockatlas.Tx{
		ID:       "30d404cc7a42b0158b95f6adfbf9a517627d60f6c7e497c1442dfdb6460285df",
		Coin:     coin.EGLD,
		Date:     int64(1587715632),
		From:     "metachain",
		To:       "erd10yagg2vme2jns9zqf9xn8kl86fkc6dr063vnuj0mz2kk2jw0qwuqmfmaw0",
		Fee:      "1000",
		Status:   blockatlas.StatusCompleted,
		Memo:     "ok",
		Sequence: 0,
		Meta: blockatlas.Transfer{
			Value:    "82516976060558456822",
			Symbol:   coin.Elrond().Symbol,
			Decimals: coin.Elrond().Decimals,
		},
		Direction: blockatlas.DirectionOutgoing,
	}

	txTransfer2Normalized = blockatlas.Tx{
		ID:       "30d404cc7a42b0158b95f6adfbf9a517627d60f6c7e497c1442dfdb6460285df",
		Coin:     coin.EGLD,
		Date:     int64(1588757256),
		From:     "erd10yagg2vme2jns9zqf9xn8kl86fkc6dr063vnuj0mz2kk2jw0qwuqmfmaw0",
		To:       "erd10yagg2vme2jns9zqf9xn8kl86fkc6dr063vnuj0mz2kk2jw0qwuqmfmaw0",
		Fee:      "1500",
		Status:   blockatlas.StatusPending,
		Memo:     "money",
		Sequence: 1,
		Meta: blockatlas.Transfer{
			Value:    "2000",
			Symbol:   coin.Elrond().Symbol,
			Decimals: coin.Elrond().Decimals,
		},
		Direction: blockatlas.DirectionSelf,
	}

	txTransfer3Normalized = blockatlas.Tx{
		ID:       "30d404cc7a42b0158b95f6adfbf9a517627d60f6c7e497c1442dfdb6460285df",
		Coin:     coin.EGLD,
		Date:     int64(1588757256),
		From:     "erd10yagg2vme2jns9zqf9xn8kl86fkc6dr063vnuj0mz2kk2jw0qwuqmfmaw0",
		To:       "erd1v0ce6rapup6rwma5sltyv05xhp33u543nex75a7j39vsz9m6squq6mxm7y",
		Fee:      "5000",
		Status:   blockatlas.StatusError,
		Memo:     "test",
		Sequence: 19,
		Meta: blockatlas.Transfer{
			Value:    "2",
			Symbol:   coin.Elrond().Symbol,
			Decimals: coin.Elrond().Decimals,
		},
		Direction: blockatlas.DirectionOutgoing,
	}

	txTransfer4Normalized = blockatlas.Tx{
		ID:       "30d404cc7a42b0158b95f6adfbf9a517627d60f6c7e497c1442dfdb6460285df",
		Coin:     coin.EGLD,
		Date:     int64(1588757256),
		From:     "erd10yagg2vme2jns9zqf9xn8kl86fkc6dr063vnuj0mz2kk2jw0qwuqmfmaw0",
		To:       "erd1v0ce6rapup6rwma5sltyv05xhp33u543nex75a7j39vsz9m6squq6mxm7y",
		Fee:      "5000",
		Status:   blockatlas.StatusPending,
		Memo:     "test",
		Sequence: 19,
		Meta: blockatlas.Transfer{
			Value:    "2",
			Symbol:   coin.Elrond().Symbol,
			Decimals: coin.Elrond().Decimals,
		},
		Direction: blockatlas.DirectionOutgoing,
	}

	txTransfer5Normalized = blockatlas.Tx{
		ID:       "30d404cc7a42b0158b95f6adfbf9a517627d60f6c7e497c1442dfdb6460285df",
		Coin:     coin.EGLD,
		Date:     int64(1588757256),
		From:     "erd10yagg2vme2jns9zqf9xn8kl86fkc6dr063vnuj0mz2kk2jw0qwuqmfmaw0",
		To:       "erd1v0ce6rapup6rwma5sltyv05xhp33u543nex75a7j39vsz9m6squq6mxm7y",
		Fee:      "5000",
		Status:   blockatlas.StatusCompleted,
		Memo:     "test",
		Sequence: 19,
		Meta: blockatlas.Transfer{
			Value:    "2",
			Symbol:   coin.Elrond().Symbol,
			Decimals: coin.Elrond().Decimals,
		},
		Direction: blockatlas.DirectionOutgoing,
	}

	txTransfer6Normalized = blockatlas.Tx{
		ID:       "30d404cc7a42b0158b95f6adfbf9a517627d60f6c7e497c1442dfdb6460285df",
		Coin:     coin.EGLD,
		From:     "erd10yagg2vme2jns9zqf9xn8kl86fkc6dr063vnuj0mz2kk2jw0qwuqmfmaw0",
		To:       "erd1v0ce6rapup6rwma5sltyv05xhp33u543nex75a7j39vsz9m6squq6mxm7y",
		Fee:      "5000",
		Status:   blockatlas.StatusCompleted,
		Memo:     "test",
		Sequence: 25,
		Block:    620,
		Date:     1596121554,
		Meta: blockatlas.Transfer{
			Value:    "2",
			Symbol:   coin.Elrond().Symbol,
			Decimals: coin.Elrond().Decimals,
		},
		Direction: blockatlas.DirectionOutgoing,
	}
)

type test struct {
	name        string
	apiResponse string
	expected    *blockatlas.Tx
}

func TestNormalize(t *testing.T) {
	testNormalize(t, &test{
		name:        "transferSuccess",
		apiResponse: txTransferSrc1,
		expected:    &txTransfer1Normalized,
	})

	testNormalize(t, &test{
		name:        "transferPending",
		apiResponse: txTransferSrc2,
		expected:    &txTransfer2Normalized,
	})

	testNormalize(t, &test{
		name:        "transferNotExecuted",
		apiResponse: txTransferSrc3,
		expected:    &txTransfer3Normalized,
	})

	testNormalize(t, &test{
		name:        "transferPendingNewStatus",
		apiResponse: txTransferSrc4,
		expected:    &txTransfer4Normalized,
	})

	testNormalize(t, &test{
		name:        "transferSuccessNewStatus",
		apiResponse: txTransferSrc5,
		expected:    &txTransfer5Normalized,
	})
}

func TestNormalizeTxs(t *testing.T) {
	var tx1, tx2, tx3 Transaction

	_ = json.Unmarshal([]byte(txTransferSrc1), &tx1)
	_ = json.Unmarshal([]byte(txTransferSrc1), &tx2)
	_ = json.Unmarshal([]byte(txTransferSrc1), &tx3)

	txs := []Transaction{tx1, tx2, tx3}
	normalizedTxs := NormalizeTxs(txs, userAddress, Block{})
	require.Equal(t, len(txs), len(normalizedTxs))
}

func testNormalize(t *testing.T, _test *test) {
	var tx Transaction
	err := json.Unmarshal([]byte(_test.apiResponse), &tx)
	if err != nil {
		t.Error(err)
		return
	}

	normalizedTx, ok := NormalizeTx(tx, tx.Sender, Block{})
	require.True(t, ok, _test.name+": cannot normalize tx")

	resJSON, err := json.Marshal(&normalizedTx)
	require.Nil(t, err)

	dstJSON, err := json.Marshal(&_test.expected)
	require.Nil(t, err)

	require.Equal(t, string(dstJSON), string(resJSON))
}

func TestNormalizeTxsFromHyperblock(t *testing.T) {
	var tx Transaction

	_ = json.Unmarshal([]byte(txTransferSrc6), &tx)
	txs := []Transaction{tx}

	normalizedTxs := NormalizeTxs(txs, userAddress, Block{
		Nonce: 620,
		Round: 659,
	})
	require.Equal(t, len(txs), len(normalizedTxs))

	require.Equal(t, []blockatlas.Tx{txTransfer6Normalized}, normalizedTxs)
}
