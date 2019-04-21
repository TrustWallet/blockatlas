package binance

import (
	"encoding/json"
	"fmt"
)

type Account struct {
	AccountNumber int       `json:"account_number"`
	Address       string    `json:"address"`
	Balances      []Balance `json:"balances"`
	PublicKey     []byte    `json:"public_key"`
	Sequence      uint64    `json:"sequence"`
}

type Balance struct {
	Symbol string `json:"symbol"`
	Free   uint64 `json:"free"`
	Locked uint64 `json:"locked"`
	Frozen uint64 `json:"frozen"`
}

type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

type Tx struct {
	BlockHeight   uint64        `json:"blockHeight"`
	Code          int           `json:"code"`
	ConfirmBlocks int           `json:"confirmBlocks"`
	Data          string        `json:"data"`
	FromAddr      string        `json:"fromAddr"`
	OrderID       string        `json:"orderId"`
	Timestamp     int64         `json:"timeStamp"`
	ToAddr        string        `json:"toAddr"`
	Age           int64         `json:"txAge"`
	MappedAsset   string        `json:"mappedTxAsset"`
	Asset         string        `json:"txAsset"`
	Fee           json.Number   `json:"txFee"`
	Hash          string        `json:"txHash"`
	Value         json.Number   `json:"value"`
}

type TxPage struct {
	Nums int  `json:"txNums"`
	Txs  []Tx `json:"txArray"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
