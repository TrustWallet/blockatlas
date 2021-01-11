package vechain

import (
	"errors"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/golibs/address"
	"github.com/trustwallet/golibs/numbers"
)

func (p *Platform) GetTokenTxsByAddress(address, token string) (blockatlas.TxPage, error) {
	if token != gasTokenAddress {
		return nil, nil
	}
	curBlock, err := p.CurrentBlockNumber()
	if err != nil {
		return nil, err
	}
	events, err := p.client.GetLogsEvent(address, token, curBlock)
	if err != nil {
		return nil, err
	}
	eventsIDs := make([]string, 0)
	for _, event := range events {
		eventsIDs = append(eventsIDs, event.Meta.TxId)
	}

	cTxs := p.getTransactionsByIDs(eventsIDs)
	txs := make(blockatlas.TxPage, 0)
	for t := range cTxs {
		txs = append(txs, t...)
	}
	return txs, nil
}

func (p *Platform) getTransactionsByIDs(ids []string) chan blockatlas.TxPage {
	txChan := make(chan blockatlas.TxPage, len(ids))
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func(i string, c chan blockatlas.TxPage) {
			defer wg.Done()
			err := p.getTransactionChannel(i, c)
			if err != nil {
				log.Error(err)
			}
		}(id, txChan)
	}
	wg.Wait()
	close(txChan)
	return txChan
}

func (p *Platform) getTransactionChannel(id string, txChan chan blockatlas.TxPage) error {
	srcTx, err := p.client.GetTransactionByID(id)
	if err != nil {
		return err
	}

	receipt, err := p.client.GetTransactionReceiptByID(id)
	if err != nil {
		return err
	}

	txs, err := p.NormalizeTokenTransaction(srcTx, receipt)
	if err != nil {
		return err
	}
	txChan <- txs
	return nil
}

func (p *Platform) NormalizeTokenTransaction(srcTx Tx, receipt TxReceipt) (blockatlas.TxPage, error) {
	txs := make(blockatlas.TxPage, 0)

	fee, err := numbers.HexToDecimal(receipt.Paid)
	if err != nil {
		return txs, err
	}

	originSender, err := address.EIP55Checksum(blockatlas.GetValidParameter(srcTx.Origin, srcTx.Meta.TxOrigin))
	if err != nil {
		return txs, err
	}

	if receipt.Reverted {
		var to string
		if len(srcTx.Clauses) > 0 {
			to = srcTx.Clauses[0].To
			if checksumTo, err := address.EIP55Checksum(to); err == nil {
				to = checksumTo
			}
		} else {
			return txs, errors.New("NormalizeBlockTransaction: srcTx.Clauses not found: " + srcTx.Id)
		}

		txs = append(txs, blockatlas.Tx{
			ID:     srcTx.Id,
			Coin:   p.Coin().ID,
			From:   originSender,
			To:     to,
			Fee:    blockatlas.Amount(fee),
			Date:   srcTx.Meta.BlockTimestamp,
			Type:   blockatlas.TxTokenTransfer,
			Block:  srcTx.Meta.BlockNumber,
			Status: blockatlas.StatusError,
		})
		return txs, nil
	}

	if receipt.Outputs == nil || len(receipt.Outputs) == 0 {
		return blockatlas.TxPage{}, errors.New("NormalizeBlockTransaction: receipt.Outputs not found: " + srcTx.Id)
	}

	for _, output := range receipt.Outputs {
		if len(output.Events) == 0 || len(output.Events[0].Topics) < 3 {
			continue
		}
		event := output.Events[0]
		value, err := numbers.HexToDecimal(event.Data)
		if err != nil {
			continue
		}

		originReceiver, err := address.EIP55Checksum(event.Address)
		if err != nil {
			continue
		}
		topicsFrom, err := address.EIP55Checksum(getRecipientAddress(event.Topics[1]))
		if err != nil {
			continue
		}
		topicsTo, err := address.EIP55Checksum(getRecipientAddress(event.Topics[2]))
		if err != nil {
			continue
		}

		direction, err := getTokenTransactionDirectory(originSender, topicsFrom, topicsTo)
		if err != nil {
			continue
		}

		txs = append(txs, blockatlas.Tx{
			ID:        srcTx.Id,
			Coin:      p.Coin().ID,
			From:      originSender,
			To:        originReceiver,
			Fee:       blockatlas.Amount(fee),
			Date:      srcTx.Meta.BlockTimestamp,
			Type:      blockatlas.TxTokenTransfer,
			Block:     srcTx.Meta.BlockNumber,
			Status:    blockatlas.StatusCompleted,
			Direction: direction,
			Meta: blockatlas.TokenTransfer{
				// the only supported Token on VeChain is its Gas token
				Name:     gasTokenName,
				TokenID:  originReceiver,
				Value:    blockatlas.Amount(value),
				Symbol:   gasTokenSymbol,
				Decimals: gasTokenDecimals,
				From:     originSender,
				To:       topicsTo,
			},
		})
	}
	return txs, nil
}

func (p *Platform) GetTxsByAddress(address string) (blockatlas.TxPage, error) {
	headBlock, err := p.CurrentBlockNumber()
	if err != nil {
		return nil, err
	}
	transfers, err := p.client.GetTransactions(address, headBlock)
	if err != nil {
		return nil, err
	}

	txs := make(blockatlas.TxPage, 0)
	for _, t := range transfers {
		trxId, err := p.client.GetTransactionByID(t.Meta.TxId)
		if err != nil {
			continue
		}
		tx, err := p.NormalizeTransaction(t, trxId, address)
		if err != nil {
			continue
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

func (p *Platform) NormalizeTransaction(srcTx LogTransfer, trxId Tx, addr string) (tx blockatlas.Tx, err error) {
	value, err := numbers.HexToDecimal(srcTx.Amount)
	if err != nil {
		return
	}

	fee := strconv.Itoa(trxId.Gas)
	sender, err := address.EIP55Checksum(srcTx.Sender)
	if err != nil {
		return
	}
	recipient, err := address.EIP55Checksum(srcTx.Recipient)
	if err != nil {
		return
	}
	addrChecksum, err := address.EIP55Checksum(addr)
	if err != nil {
		return
	}

	directory, err := getTransferDirectory(sender, recipient, addrChecksum)
	if err != nil {
		return blockatlas.Tx{}, err
	}

	return blockatlas.Tx{
		ID:        srcTx.Meta.TxId,
		Coin:      p.Coin().ID,
		From:      sender,
		To:        recipient,
		Fee:       blockatlas.Amount(fee),
		Date:      srcTx.Meta.BlockTimestamp,
		Type:      blockatlas.TxTransfer,
		Block:     srcTx.Meta.BlockNumber,
		Direction: directory,
		Status:    blockatlas.StatusCompleted,
		Meta: blockatlas.Transfer{
			Value:    blockatlas.Amount(value),
			Symbol:   p.Coin().Symbol,
			Decimals: p.Coin().Decimals,
		},
	}, nil
}

func hexToInt(hex string) (uint64, error) {
	nonceStr, err := numbers.HexToDecimal(hex)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(nonceStr, 10, 64)
}

// Substring recipient address from clause data and appends 0x
// 0x000000000000000000000000b5e883349e68ab59307d1604555ac890fac47128 => 0xb5e883349e68ab59307d1604555ac890fac47128
func getRecipientAddress(hex string) string {
	return "0x" + hex[len(hex)-40:]
}

func getTokenTransactionDirectory(originSender, topicsFrom, topicsTo string) (dir blockatlas.Direction, err error) {
	if originSender == topicsFrom && originSender == topicsTo {
		return blockatlas.DirectionSelf, nil
	}
	if originSender == topicsFrom && originSender != topicsTo {
		return blockatlas.DirectionIncoming, nil
	}
	if originSender == topicsTo && originSender != topicsFrom {
		return blockatlas.DirectionOutgoing, nil
	}
	return "", errors.New("Unknown direction")
}

func getTransferDirectory(sender, recipient, addr string) (dir blockatlas.Direction, err error) {
	if sender == addr && recipient == addr {
		return blockatlas.DirectionSelf, nil
	}
	if sender == addr && recipient != addr {
		return blockatlas.DirectionOutgoing, nil
	}
	if recipient == addr && sender != addr {
		return blockatlas.DirectionIncoming, nil
	}
	return "", errors.New("Unknown direction")
}
