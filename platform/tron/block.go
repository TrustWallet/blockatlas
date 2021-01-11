package tron

import (
	"encoding/hex"
	"sync"

	"github.com/trustwallet/golibs/txtype"
)

func (p *Platform) CurrentBlockNumber() (int64, error) {
	return p.client.fetchCurrentBlockNumber()
}

func (p *Platform) GetBlockByNumber(num int64) (*txtype.Block, error) {
	block, err := p.client.fetchBlockByNumber(num)
	if err != nil {
		return nil, err
	}

	txsChan := p.NormalizeBlockTxs(block.Txs)
	txs := make(txtype.TxPage, 0)
	for cTxs := range txsChan {
		txs = append(txs, cTxs)
	}

	return &txtype.Block{
		Number: num,
		Txs:    txs,
	}, nil
}

func (p *Platform) NormalizeBlockTxs(srcTxs []Tx) chan txtype.Tx {
	txChan := make(chan txtype.Tx, len(srcTxs))
	var wg sync.WaitGroup
	for _, srcTx := range srcTxs {
		wg.Add(1)
		go func(s Tx, c chan txtype.Tx) {
			defer wg.Done()
			p.NormalizeBlockChannel(s, c)
		}(srcTx, txChan)
	}
	wg.Wait()
	close(txChan)
	return txChan
}

func (p *Platform) NormalizeBlockChannel(srcTx Tx, txChan chan txtype.Tx) {
	if len(srcTx.Data.Contracts) == 0 {
		return
	}

	tx, err := normalize(srcTx)
	if err != nil {
		return
	}
	transfer := srcTx.Data.Contracts[0].Parameter.Value
	if len(transfer.AssetName) > 0 {
		assetName, err := hex.DecodeString(transfer.AssetName[:])
		if err == nil {
			info, err := p.client.fetchTokenInfo(string(assetName))
			if err == nil && len(info.Data) > 0 {
				addTokenMeta(tx, srcTx, info.Data[0])
			}
		}
	}
	txChan <- *tx
}
