package tron

import (
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/address"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"strconv"
)

func (p *Platform) GetTxsByAddress(address string) (blockatlas.TxPage, error) {
	Txs, err := p.client.fetchTxsOfAddress(address, "")
	if err != nil && len(Txs) == 0 {
		return nil, err
	}

	txs := make(blockatlas.TxPage, 0)
	for _, srcTx := range Txs {
		tx, err := normalize(srcTx)
		if err != nil {
			continue
		}

		if len(srcTx.Data.Contracts) > 0 && srcTx.Data.Contracts[0].Type == TransferContract {
			txs = append(txs, *tx)
		} else {
			continue
		}
	}

	return txs, nil
}

func (p *Platform) GetTokenTxsByAddress(address, token string) (blockatlas.TxPage, error) {
	unknownTokenType := errors.E("unknownTokenType")
	tokenType := getTokenType(token)

	switch tokenType {
	case blockatlas.TokenTypeTRC10:
		txs, err := p.fetchTransactionsForTRC10Tokens(address, token)
		if err != nil {
			return nil, err
		}
		return txs, nil
	case blockatlas.TokenTypeTRC20:
		trc20Transactions, err := p.client.fetchTRC20Transactions(address)
		if err != nil {
			return nil, err
		}
		return blockatlas.TxPage(normalizeTRC20Transactions(trc20Transactions)), nil
	default:
		return nil, unknownTokenType
	}
}

func getTokenType(token string) blockatlas.TokenType {
	_, err := strconv.Atoi(token)
	if err != nil {
		return blockatlas.TokenTypeTRC20
	} else {
		return blockatlas.TokenTypeTRC10
	}
}

func addTokenMeta(tx *blockatlas.Tx, srcTx Tx, tokenInfo AssetInfo) {
	transfer := srcTx.Data.Contracts[0].Parameter.Value
	tx.Meta = blockatlas.TokenTransfer{
		Name:     tokenInfo.Name,
		Symbol:   tokenInfo.Symbol,
		TokenID:  tokenInfo.ID,
		Decimals: tokenInfo.Decimals,
		Value:    transfer.Amount,
		From:     tx.From,
		To:       tx.To,
	}
}

func (p *Platform) fetchTransactionsForTRC10Tokens(address, token string) (blockatlas.TxPage, error) {
	txs := make(blockatlas.TxPage, 0)

	tokenTxs, err := p.client.fetchTxsOfAddress(address, token)
	if err != nil {
		return nil, errors.E(err, "TRON: failed to get token from address", errors.TypePlatformApi,
			errors.Params{"address": address, "token": token})
	}

	info, err := p.client.fetchTokenInfo(token)
	if err != nil {
		return nil, errors.E(err, "TRON: failed to get token info", errors.TypePlatformApi,
			errors.Params{"address": address, "token": token})
	}
	for _, srcTx := range tokenTxs {
		tx, err := normalize(srcTx)
		if err != nil {
			logger.Error(err)
			continue
		}
		if info.Data != nil && len(info.Data) > 0 {
			addTokenMeta(tx, srcTx, info.Data[0])
		}

		txs = append(txs, *tx)
	}
	return txs, nil
}

func normalizeTRC20Transactions(transactions TRC20Transactions) blockatlas.Txs {
	txs := make(blockatlas.Txs, 0, len(transactions.Data))
	for _, rawTx := range transactions.Data {
		tx := blockatlas.Tx{
			ID:     rawTx.TransactionID,
			Coin:   coin.TRX,
			Date:   rawTx.BlockTimestamp / 1000,
			From:   rawTx.From,
			To:     rawTx.To,
			Fee:    "0",
			Block:  0,
			Status: blockatlas.StatusCompleted,
			Meta: blockatlas.TokenTransfer{
				Name:     rawTx.TokenInfo.Name,
				Symbol:   rawTx.TokenInfo.Symbol,
				TokenID:  rawTx.TokenInfo.Address,
				Decimals: uint(rawTx.TokenInfo.Decimals),
				Value:    blockatlas.Amount(rawTx.Value),
				From:     rawTx.From,
				To:       rawTx.To,
			},
		}
		txs = append(txs, tx)
	}
	return txs
}

func normalize(srcTx Tx) (*blockatlas.Tx, error) {
	if len(srcTx.Data.Contracts) == 0 {
		return nil, errors.E("TRON: transfer without contract", errors.TypePlatformApi,
			errors.Params{"tx": srcTx})
	}

	contract := srcTx.Data.Contracts[0]
	if contract.Type != TransferContract && contract.Type != TransferAssetContract {
		return nil, errors.E("TRON: invalid contract transfer", errors.TypePlatformApi,
			errors.Params{"tx": srcTx, "type": contract.Type})
	}

	transfer := contract.Parameter.Value
	from, err := address.HexToAddress(transfer.OwnerAddress)
	if err != nil {
		return nil, errors.E(err, "TRON: failed to get from address", errors.TypePlatformApi,
			errors.Params{"tx": srcTx})
	}
	to, err := address.HexToAddress(transfer.ToAddress)
	if err != nil {
		return nil, errors.E(err, "TRON: failed to get to address", errors.TypePlatformApi,
			errors.Params{"tx": srcTx})
	}

	return &blockatlas.Tx{
		ID:     srcTx.ID,
		Coin:   coin.TRX,
		Date:   srcTx.BlockTime / 1000,
		From:   from,
		To:     to,
		Fee:    "0",
		Block:  0,
		Status: blockatlas.StatusCompleted,
		Meta: blockatlas.Transfer{
			Value:    transfer.Amount,
			Symbol:   coin.Coins[coin.TRX].Symbol,
			Decimals: coin.Coins[coin.TRX].Decimals,
		},
	}, nil
}
