package models

import (
	"encoding/json"
	"fmt"
)

// Tx, but with default JSON marshalling methods
type wrappedTx Tx

func (t *Tx) UnmarshalJSON(data []byte) error {
	// Wrap the Tx type to avoid infinite recursion
	var wrapped wrappedTx

	var raw json.RawMessage
	wrapped.Meta = &raw
	if err := json.Unmarshal(data, &wrapped); err != nil {
		return err
	}

	*t = Tx(wrapped)

	switch t.Type {
	case TxTransfer:
		t.Meta = new(Transfer)
	case TxTokenTransfer:
		t.Meta = new(TokenTransfer)
	case TxCollectibleTransfer:
		t.Meta = new(CollectibleTransfer)
	case TxTokenSwap:
		t.Meta = new(TokenSwap)
	case TxContractCall:
		t.Meta = new(ContractCall)
	default:
		return fmt.Errorf(`unsupported tx type "%s"`, t.Type)
	}
	if err := json.Unmarshal(raw, t.Meta); err != nil {
		return err
	}
	return nil
}

func (t *Tx) MarshalJSON() ([]byte, error) {
	// Set type from metadata content
	switch t.Meta.(type) {
	case Transfer, *Transfer:
		t.Type = TxTransfer
	case TokenTransfer, *TokenTransfer:
		t.Type = TxTokenTransfer
	case CollectibleTransfer, *CollectibleTransfer:
		t.Type = TxCollectibleTransfer
	case TokenSwap, *TokenSwap:
		t.Type = TxTokenSwap
	case ContractCall, *ContractCall:
		t.Type = TxContractCall
	default:
		return nil, fmt.Errorf("unsupported tx metadata")
	}

	// Wrap the Tx type to avoid infinite recursion
	return json.Marshal(wrappedTx(*t))
}
