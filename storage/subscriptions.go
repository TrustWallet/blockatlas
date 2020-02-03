package storage

import (
	"fmt"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"sync"
)

const (
	ATLAS_OBSERVER = "ATLAS_OBSERVER"
)

func (s *Storage) Lookup(coin uint, addresses []string) ([]blockatlas.Subscription, error) {
	if len(addresses) == 0 {
		return nil, errors.E("cannot look up an empty list")
	}

	observers := make([]blockatlas.Subscription, 0)
	var wg sync.WaitGroup
	out := make(chan []blockatlas.Subscription)
	wg.Add(len(addresses))
	for _, address := range addresses {
		go func(coin uint, addr string) {
			defer wg.Done()
			key := getSubscriptionKey(coin, addr)
			var webhooks []string
			err := s.GetHMValue(ATLAS_OBSERVER, key, &webhooks)
			if err != nil {
				return
			}
			subs := make([]blockatlas.Subscription, 0)
			for _, webhook := range webhooks {
				subs = append(subs, blockatlas.Subscription{Coin: coin, Address: addr, Webhook: webhook})
			}
			out <- subs
		}(coin, address)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	for r := range out {
		observers = append(observers, r...)
	}
	return observers, nil
}

func (s *Storage) AddSubscriptions(subscriptions []blockatlas.Subscription) {
	for _, sub := range subscriptions {
		key := getSubscriptionKey(sub.Coin, sub.Address)
		var webhooks []string
		_ = s.GetHMValue(ATLAS_OBSERVER, key, &webhooks)
		if webhooks == nil {
			webhooks = make([]string, 0)
		}
		if hasObject(webhooks, sub.Webhook) {
			continue
		}
		webhooks = append(webhooks, sub.Webhook)
		err := s.AddHM(ATLAS_OBSERVER, key, webhooks)
		if err != nil {
			logger.Error(err, "AddSubscriptions error", errors.Params{"webhooks": webhooks, "address": sub.Address, "coin": sub.Coin})
		}
	}
}

func (s *Storage) DeleteSubscriptions(subscriptions []blockatlas.Subscription) {
	for _, sub := range subscriptions {
		key := getSubscriptionKey(sub.Coin, sub.Address)
		var webhooks []string
		err := s.GetHMValue(ATLAS_OBSERVER, key, &webhooks)
		if err != nil {
			continue
		}
		newHooks := make([]string, 0)
		for _, webhook := range webhooks {
			if webhook == sub.Webhook {
				continue
			}
			newHooks = append(newHooks, webhook)
		}
		if len(newHooks) == 0 {
			_ = s.DeleteHM(ATLAS_OBSERVER, key)
			continue
		}
		err = s.AddHM(ATLAS_OBSERVER, key, newHooks)
		if err != nil {
			logger.Error(err, "DeleteSubscriptions - AddHM", errors.Params{"webhook": newHooks, "address": sub.Address, "coin": sub.Coin})
		}
	}
}

func getSubscriptionKey(coin uint, address string) string {
	return fmt.Sprintf("%d-%s", coin, address)
}

func hasObject(array []string, obj string) bool {
	for _, temp := range array {
		if temp == obj {
			return true
		}
	}
	return false
}
