package iotex

import (
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"strconv"
	"time"

	"github.com/trustwallet/golibs/coin"
)

func (p *Platform) GetTxsByAddress(address string) (blockatlas.TxPage, error) {
	txs := make([]blockatlas.Tx, 0)
	var start int64

	totalTrx, err := p.client.GetAddressTotalTransactions(address)
	if err != nil {
		return nil, err
	}

	if totalTrx >= blockatlas.TxPerPage {
		start = totalTrx - blockatlas.TxPerPage
	}

	actions, err := p.client.GetTxsOfAddress(address, start)
	if err != nil {
		return nil, err
	}

	for _, srcTx := range actions.ActionInfo {
		tx := Normalize(srcTx)
		if tx != nil {
			txs = append(txs, *tx)
		}
	}

	return txs, nil
}

// Normalize converts an Iotex transaction into the generic model
func Normalize(trx *ActionInfo) *blockatlas.Tx {
	if trx.Action == nil {
		return nil
	}
	if trx.Action.Core == nil {
		return nil
	}
	if trx.Action.Core.Transfer == nil {
		return nil
	}

	date, err := time.Parse(time.RFC3339, trx.Timestamp)
	if err != nil {
		return nil
	}
	height, err := strconv.ParseInt(trx.BlkHeight, 10, 64)
	if err != nil {
		return nil
	}
	if height <= 0 {
		return nil
	}
	nonce, err := strconv.ParseInt(trx.Action.Core.Nonce, 10, 64)
	if err != nil {
		return nil
	}

	return &blockatlas.Tx{
		ID:       trx.ActHash,
		Coin:     coin.IOTX,
		From:     trx.Sender,
		To:       trx.Action.Core.Transfer.Recipient,
		Fee:      blockatlas.Amount(trx.GasFee),
		Date:     date.Unix(),
		Block:    uint64(height),
		Status:   blockatlas.StatusCompleted,
		Sequence: uint64(nonce),
		Type:     blockatlas.TxTransfer,
		Meta: blockatlas.Transfer{
			Value:    trx.Action.Core.Transfer.Amount,
			Symbol:   coin.Coins[coin.IOTX].Symbol,
			Decimals: coin.Coins[coin.IOTX].Decimals,
		},
	}
}
