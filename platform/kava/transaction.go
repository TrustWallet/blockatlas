package kava

import (
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/golibs/numbers"
)

const KAVADenom = "ukava"

func (p *Platform) GetTxsByAddress(address string) (blockatlas.TxPage, error) {
	return p.GetTokenTxsByAddress(address, KAVADenom)
}

func (p *Platform) GetTokenTxsByAddress(address, token string) (blockatlas.TxPage, error) {
	tagsList := []string{"transfer.recipient", "message.sender"}
	var wg sync.WaitGroup
	out := make(chan []Tx, len(tagsList))
	wg.Add(len(tagsList))
	for _, t := range tagsList {
		go func(tag, addr string, wg *sync.WaitGroup) {
			defer wg.Done()
			page := 1
			txs, err := p.client.GetAddrTxs(addr, tag, page)
			if err != nil {
				log.WithFields(log.Fields{"address": tag, "tag": tag}).Error("GetAddrTxs", err)
				return
			}
			// Condition when no more pages to paginate
			if txs.PageTotal == "1" {
				out <- txs.Txs
				return
			}

			totalPages, err := strconv.Atoi(txs.PageTotal)
			if err != nil {
				log.WithFields(log.Fields{"totalPages": totalPages}).Error("GetAddrTxs", err)
				return
			}
			// gaia does support sort option, paginate to get latest transactions by passing total pages page
			// https://github.com/cosmos/gaia/blob/f61b391aee5d04364d2b5539692bbb187ad9b946/docs/resources/gaiacli.md#query-transactions
			txs2, err := p.client.GetAddrTxs(addr, tag, totalPages)
			if err != nil {
				log.WithFields(log.Fields{"address": tag, "tag": tag}).Error("GetAddrTxs", err)
				return
			}
			out <- txs2.Txs
		}(t, address, &wg)
	}
	wg.Wait()
	close(out)
	srcTxs := make([]Tx, 0)
	for r := range out {
		filteredTxs := p.FilterTxsByDenom(r, token)
		srcTxs = append(srcTxs, filteredTxs...)
	}
	return p.NormalizeTxs(srcTxs), nil
}

func (p *Platform) FilterTxsByDenom(txs []Tx, denom string) []Tx {
	filteredTxs := make([]Tx, 0)
	for _, tx := range txs {
		if tx.Data.Contents.Message[0].Value.(MessageValueTransfer).Amount[0].Denom == denom {
			filteredTxs = append(filteredTxs, tx)
		}
	}
	return filteredTxs
}

// NormalizeTxs converts multiple Cosmos transactions
func (p *Platform) NormalizeTxs(srcTxs []Tx) blockatlas.TxPage {
	txMap := make(map[string]bool)
	txs := make(blockatlas.TxPage, 0)
	for _, srcTx := range srcTxs {
		_, ok := txMap[srcTx.ID]
		if ok {
			continue
		}
		normalisedInputTx, ok := p.Normalize(&srcTx)
		if ok {
			txMap[srcTx.ID] = true
			txs = append(txs, normalisedInputTx)
		}
	}
	return txs
}

// Normalize converts an Cosmos transaction into the generic model
func (p *Platform) Normalize(srcTx *Tx) (tx blockatlas.Tx, ok bool) {
	date, err := time.Parse("2006-01-02T15:04:05Z", srcTx.Date)
	if err != nil {
		return blockatlas.Tx{}, false
	}
	block, err := strconv.ParseUint(srcTx.Block, 10, 64)
	if err != nil {
		return blockatlas.Tx{}, false
	}
	// Sometimes fees can be null objects (in the case of no fees e.g. F044F91441C460EDCD90E0063A65356676B7B20684D94C731CF4FAB204035B41)
	fee := "0"
	if len(srcTx.Data.Contents.Fee.FeeAmount) > 0 {
		qty := srcTx.Data.Contents.Fee.FeeAmount[0].Quantity
		if len(qty) > 0 && qty != fee {
			fee, err = numbers.DecimalToSatoshis(srcTx.Data.Contents.Fee.FeeAmount[0].Quantity)
			if err != nil {
				return blockatlas.Tx{}, false
			}
		}
	}

	status := blockatlas.StatusCompleted
	// https://github.com/cosmos/cosmos-sdk/blob/95ddc242ad024ca78a359a13122dade6f14fd676/types/errors/errors.go#L19
	if srcTx.Code > 0 {
		status = blockatlas.StatusError
	}

	tx = blockatlas.Tx{
		ID:     srcTx.ID,
		Coin:   p.Coin().ID,
		Date:   date.Unix(),
		Status: status,
		Fee:    blockatlas.Amount(fee),
		Block:  block,
		Memo:   srcTx.Data.Contents.Memo,
	}

	if len(srcTx.Data.Contents.Message) == 0 {
		return tx, false
	}

	msg := srcTx.Data.Contents.Message[0]
	switch msg.Value.(type) {
	case MessageValueTransfer:
		transfer := msg.Value.(MessageValueTransfer)
		p.fillTransfer(&tx, transfer)
		return tx, true
	case MessageValueDelegate:
		delegate := msg.Value.(MessageValueDelegate)
		p.fillDelegate(&tx, delegate, srcTx.Events, msg.Type)
		return tx, true
	}
	return tx, false
}

func (p *Platform) fillTransfer(tx *blockatlas.Tx, transfer MessageValueTransfer) {
	if len(transfer.Amount) == 0 {
		return
	}
	value, err := numbers.DecimalToSatoshis(transfer.Amount[0].Quantity)
	if err != nil {
		return
	}
	tx.From = transfer.FromAddr
	tx.To = transfer.ToAddr
	tx.Type = blockatlas.TxTransfer
	tx.Meta = blockatlas.Transfer{
		Value:    blockatlas.Amount(value),
		Symbol:   p.Coin().Symbol,
		Decimals: p.Coin().Decimals,
	}
	switch {
	case transfer.Amount[0].Denom == KAVADenom:
		tx.Type = blockatlas.TxTransfer
		tx.Meta = blockatlas.Transfer{
			Value:    blockatlas.Amount(value),
			Symbol:   p.Coin().Symbol,
			Decimals: p.Coin().Decimals,
		}
	default:
		tx.Type = blockatlas.TxNativeTokenTransfer
		tx.Meta = blockatlas.NativeTokenTransfer{
			Decimals: p.Coin().Decimals,
			From:     tx.From,
			Symbol:   strings.ToUpper(transfer.Amount[0].Denom),
			Name:     transfer.Amount[0].Denom,
			To:       tx.To,
			TokenID:  transfer.Amount[0].Denom,
			Value:    blockatlas.Amount(value),
		}
	}
}

func (p *Platform) fillDelegate(tx *blockatlas.Tx, delegate MessageValueDelegate, events Events, msgType TxType) {
	value := ""
	if len(delegate.Amount.Quantity) > 0 {
		var err error
		value, err = numbers.DecimalToSatoshis(delegate.Amount.Quantity)
		if err != nil {
			return
		}
	}
	tx.From = delegate.DelegatorAddr
	tx.To = delegate.ValidatorAddr
	tx.Type = blockatlas.TxAnyAction

	key := blockatlas.KeyStakeDelegate
	title := blockatlas.KeyTitle("")
	switch msgType {
	case MsgDelegate:
		tx.Direction = blockatlas.DirectionOutgoing
		title = blockatlas.AnyActionDelegation
	case MsgUndelegate:
		tx.Direction = blockatlas.DirectionIncoming
		title = blockatlas.AnyActionUndelegation
	case MsgWithdrawDelegationReward:
		tx.Direction = blockatlas.DirectionIncoming
		title = blockatlas.AnyActionClaimRewards
		key = blockatlas.KeyStakeClaimRewards
		value = events.GetWithdrawRewardValue()
	}
	tx.Meta = blockatlas.AnyAction{
		Coin:     p.Coin().ID,
		Title:    title,
		Key:      key,
		Name:     p.Coin().Name,
		Symbol:   p.Coin().Symbol,
		Decimals: p.Coin().Decimals,
		Value:    blockatlas.Amount(value),
	}
}
