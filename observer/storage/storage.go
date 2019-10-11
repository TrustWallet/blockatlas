package storage

import (
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/blockatlas/pkg/storage/sql"
	"github.com/trustwallet/blockatlas/platform/bitcoin"
)

type Storage struct {
	sql.PgSql
}

func (s *Storage) GetBlockNumber(coin uint) (int64, error) {
	b := Block{Coin: coin}
	err := s.Get(&b)
	if err != nil {
		return 0, nil
	}
	return b.BlockHeight, nil
}

func (s *Storage) SetBlockNumber(coin uint, num int64) error {
	b := Block{Coin: coin, BlockHeight: num}
	err := s.CreateOrUpdate(&b)
	if err != nil {
		return errors.E(err, errors.Params{"block": num, "coin": coin}).PushToSentry()
	}
	return nil
}

func (s *Storage) SaveXpubAddresses(coin uint, addresses []string, xpub string) error {
	if len(addresses) == 0 {
		return errors.E("no addresses for xpub", errors.Params{"xpub": xpub}).PushToSentry()
	}

	a := make([]interface{}, 0)
	for _, address := range addresses {
		x := &Xpub{
			Xpub:    xpub,
			Address: address,
			Coin:    coin,
		}
		a = append(a, x)
	}
	return s.AddMany(a...)
}

func (s *Storage) GetAddressFromXpub(coin uint, xpub string) ([]Xpub, error) {
	x := &Xpub{
		Xpub: xpub,
		Coin: coin,
	}

	var addresses []Xpub
	err := s.Find(&addresses, &x)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (s *Storage) GetXpubFromAddress(coin uint, address string) (string, error) {
	a := &Xpub{
		Address: address,
	}
	err := s.Get(&a)
	if err != nil {
		return "", err
	}
	return a.Xpub, nil
}

func (s *Storage) Lookup(coin uint, addresses ...string) (observers []Subscription, err error) {
	if len(addresses) == 0 {
		return nil, errors.E("cannot look up an empty list", errors.Params{"coin": coin}).PushToSentry()
	}
	s.Client.
		Table("subscriptions").
		Select("subscriptions.coin, COALESCE(xpubs.address, subscriptions.address) AS address, subscriptions.webhook, xpubs.xpub AS xpub").
		Joins("LEFT JOIN xpubs ON subscriptions.address = xpubs.xpub").
		Where("subscriptions.address IN (?)", addresses).
		Or("xpubs.address IN (?)", addresses).
		Find(&observers)
	return
}

func (s *Storage) AddSubscriptions(subscriptions []interface{}) error {
	return s.AddMany(subscriptions...)
}

func (s *Storage) DeleteSubscriptions(subscriptions []interface{}) error {
	return s.DeleteMany(subscriptions...)
}

func (s *Storage) CacheXPubAddress(xpub string, coin uint) {
	platform := bitcoin.UtxoPlatform(coin)
	addresses, err := platform.GetAddressesFromXpub(xpub)
	if err != nil || len(addresses) == 0 {
		logger.Error("GetAddressesFromXpub", err, logger.Params{
			"xpub":      xpub,
			"coin":      coin,
			"addresses": addresses,
		})
		return
	}
	err = s.SaveXpubAddresses(coin, addresses, xpub)
	if err != nil {
		logger.Error("SaveXpubAddresses", err, logger.Params{
			"xpub": xpub,
			"coin": coin,
		})
	}
}
