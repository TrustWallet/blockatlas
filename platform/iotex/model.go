package iotex

import "github.com/trustwallet/blockatlas"

type Response struct {
	ActionInfo []*ActionInfo `json:"actionInfo"`
}

type AccountInfo struct {
	AccountMeta *AccountMeta `json:"accountMeta"`
}

type AccountMeta struct {
	NumActions   string   `json:"numActions"`
}

type ActionInfo struct {
	Action    *Action `json:"action"`
	ActHash   string  `json:"actHash"`
	BlkHeight string  `json:"blkHeight"`
	Sender    string  `json:"sender"`
	GasFee    string  `json:"gasFee"`
	Timestamp string  `json:"timestamp"`
}

type Action struct {
	Core         *ActionCore `json:"core"`
}

type ActionCore struct {
	Nonce    string    `json:"nonce"`
	Transfer *Transfer `json:"transfer"`
}

type Transfer struct {
	Amount    blockatlas.Amount `json:"amount"`
	Recipient string `json:"recipient"`
}
