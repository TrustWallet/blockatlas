package blockatlas

import "strconv"

type (
	Subscriptions map[string][]string

	SubscriptionOperation string

	SubscriptionEvent struct {
		NewSubscriptions Subscriptions         `json:"new_subscriptions"`
		OldSubscriptions Subscriptions         `json:"old_subscriptions"`
		GUID             string                `json:"guid"`
		Operation        SubscriptionOperation `json:"operation"`
	}

	Subscription struct {
		Coin    uint   `json:"coin"`
		Address string `json:"address"`
		GUID    string `json:"guid"`
	}

	CoinStatus struct {
		Height int64  `json:"height"`
		Error  string `json:"error,omitempty"`
	}

	Observer struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}

	Block struct {
		Number int64  `json:"number"`
		ID     string `json:"id,omitempty"`
		Txs    []Tx   `json:"txs"`
	}
)

func (e *SubscriptionEvent) ParseSubscriptions(s Subscriptions) []Subscription {
	subs := make([]Subscription, 0)
	for coinStr, perCoin := range s {
		coin, err := strconv.Atoi(coinStr)
		if err != nil {
			continue
		}
		for _, addr := range perCoin {
			subs = append(subs, Subscription{
				Coin:    uint(coin),
				Address: addr,
				GUID:    e.GUID,
			})
		}
	}
	return subs
}
