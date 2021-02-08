package near

import (
	"encoding/json"
	"errors"

	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/golibs/types"
)

func (p *Platform) GetTxsByAddress(address string) (types.TxPage, error) {
	normalized := make(types.TxPage, 0)
	return normalized, nil
}

func NormalizeChunk(chunk ChunkDetail) types.TxPage {
	normalized := make(types.TxPage, 0)
	for _, tx := range chunk.Transactions {
		if len(tx.Actions) != 1 {
			continue
		}

		transfer, err := mapTransfer(tx.Actions[0])
		if err != nil {
			continue
		}

		normalized = append(normalized, types.Tx{
			ID:       tx.Hash,
			Coin:     coin.NEAR,
			From:     tx.SignerID,
			To:       tx.ReceiverID,
			Fee:      "0",
			Date:     int64(chunk.Header.Timestamp),
			Block:    chunk.Header.Height,
			Status:   types.StatusCompleted,
			Sequence: uint64(tx.Nonce),
			Type:     types.TxTransfer,
			Meta: types.Transfer{
				Value:    types.Amount(transfer.Transfer.Deposit),
				Symbol:   coin.Near().Name,
				Decimals: coin.Near().Decimals,
			},
		})
	}

	return normalized
}

func mapTransfer(i interface{}) (action TransferAction, err error) {
	bytes, err := json.Marshal(i)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &action)
	if err != nil {
		return
	}

	if action.Transfer.Deposit == "" {
		err = errors.New("unable marshalling to transfer actoin struct")
	}
	return
}
