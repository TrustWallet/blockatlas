package polkadot

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
)

type Platform struct {
	client    Client
	CoinIndex uint
}

func (p *Platform) Init() error {
	p.client = Client{blockatlas.InitClient(viper.GetString(p.ConfigKey()))}
	return nil
}

func (p *Platform) Coin() coin.Coin {
	return coin.Coins[p.CoinIndex]
}

func (p *Platform) ConfigKey() string {
	return fmt.Sprintf("%s.api", p.Coin().Handle)
}

func (p *Platform) GetTxsByAddress(address string) (blockatlas.TxPage, error) {
	transfers, err := p.client.GetTransfersOfAddress(address)
	if err != nil {
		return nil, err
	}

	txs := make([]blockatlas.Tx, 0)
	for _, srcTx := range transfers {
		tx := p.NormalizeTransfer(&srcTx)
		txs = append(txs, tx)
	}

	return txs, nil
}

func (p *Platform) CurrentBlockNumber() (int64, error) {
	return p.client.GetCurrentBlock()
}

func (p *Platform) GetBlockByNumber(num int64) (*blockatlas.Block, error) {
	if srcBlock, err := p.client.GetBlockByNumber(num); err == nil {
		txs := p.NormalizeExtrinsics(srcBlock)
		return &blockatlas.Block{
			Number: num,
			Txs:    txs,
		}, nil
	} else {
		return nil, err
	}
}

func (p *Platform) NormalizeTransfer(srcTx *Transfer) blockatlas.Tx {
	decimals := p.Coin().Decimals
	amount := ParseAmount(srcTx.Amount, decimals)
	status := blockatlas.StatusCompleted
	if !srcTx.Success {
		status = blockatlas.StatusFailed
	}
	result := blockatlas.Tx{
		ID:     srcTx.Hash,
		Coin:   p.Coin().ID,
		Date:   int64(srcTx.Timestamp),
		From:   srcTx.From,
		To:     srcTx.To,
		Fee:    blockatlas.Amount(FeeTransfer), // API will return fee later
		Block:  srcTx.BlockNumber,
		Status: status,
		Meta: blockatlas.Transfer{
			Value:    blockatlas.Amount(amount),
			Symbol:   p.Coin().Symbol,
			Decimals: decimals,
		},
	}
	return result
}

func (p *Platform) NormalizeExtrinsics(extrinsics []Extrinsic) []blockatlas.Tx {
	txs := make([]blockatlas.Tx, 0)
	for _, srcTx := range extrinsics {
		tx := p.NormalizeExtrinsic(&srcTx)
		if tx != nil {
			txs = append(txs, *tx)
		}
	}
	return txs
}

func (p *Platform) NormalizeExtrinsic(srcTx *Extrinsic) *blockatlas.Tx {
	var datas []CallData
	err := json.Unmarshal([]byte(srcTx.Params), &datas)
	if err != nil {
		return nil
	}

	var status blockatlas.Status
	if !srcTx.Success {
		status = blockatlas.StatusFailed
	} else {
		status = blockatlas.StatusCompleted
	}

	result := blockatlas.Tx{
		ID:       srcTx.Hash,
		Coin:     p.Coin().ID,
		Date:     int64(srcTx.Timestamp),
		Block:    srcTx.BlockNumber,
		Status:   status,
		Sequence: srcTx.Nonce,
	}

	decimals := p.Coin().Decimals
	if srcTx.CallModule == ModuleBalances &&
		srcTx.CallModuleFunction == ModuleFunctionTransfer {
		result.From = srcTx.AccountId // FIXME
		result.To = datas[0].Value.(string)
		result.Fee = blockatlas.Amount(FeeTransfer)
		result.Meta = blockatlas.Transfer{
			Value:    blockatlas.Amount(fmt.Sprintf("%.0f", datas[1].Value.(float64))),
			Symbol:   p.Coin().Symbol,
			Decimals: decimals,
		}
	} else {
		// not supported yet
		return nil
	}
	return &result
}
