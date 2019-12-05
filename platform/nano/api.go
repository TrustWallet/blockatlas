package nano

import (
	"strconv"

	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"

	"github.com/spf13/viper"
)

type Platform struct {
	client Client
}

func (p *Platform) Init() error {
	p.client = Client{blockatlas.InitClient(viper.GetString("nano.api"))}
	return nil
}

func (p *Platform) Coin() coin.Coin {
	return coin.Coins[coin.NANO]
}

func (p *Platform) GetTxsByAddress(address string) (blockatlas.TxPage, error) {
	normalized := make([]blockatlas.Tx, 0)
	history, err := p.client.GetAccountHistory(address)
	if err != nil {
		return nil, err
	}

	for _, srcTx := range history.History {
		tx := p.Normalize(&srcTx, history.Account)
		normalized = append(normalized, tx)
	}

	return normalized, nil
}

func (p *Platform) Normalize(srcTx *Transaction, account string) (tx blockatlas.Tx) {
	var from string
	var to string

	if srcTx.Type == BlockTypeSend {
		from = account
		to = srcTx.Account
	} else if srcTx.Type == BlockTypeReceive {
		from = srcTx.Account
		to = account
	}

	status := blockatlas.StatusCompleted
	height, _ := strconv.ParseUint(srcTx.Height, 10, 64)
	if height == 0 {
		status = blockatlas.StatusPending
	}
	timestamp, _ := strconv.ParseInt(srcTx.LocalTimestamp, 10, 64)

	tx = blockatlas.Tx{
		ID:     srcTx.Hash,
		Coin:   p.Coin().ID,
		Date:   timestamp,
		From:   from,
		To:     to,
		Block:  height,
		Status: status,
		Meta: blockatlas.Transfer{
			Value:    blockatlas.Amount(srcTx.Amount),
			Symbol:   p.Coin().Symbol,
			Decimals: p.Coin().Decimals,
		},
	}
	return tx
}
