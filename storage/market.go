package storage

import (
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/errors"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"strings"
)

const (
	EntityRates  = "ATLAS_MARKET_RATES"
	EntityQuotes = "ATLAS_MARKET_QUOTES"
)

type MarketProviderList interface {
	GetPriority(providerId string) int
}

func (s *Storage) SaveTicker(coin blockatlas.Ticker, pl MarketProviderList) error {
	cd, err := s.GetTicker(coin.CoinName, coin.TokenId)
	if err == nil {
		op := pl.GetPriority(cd.Price.Provider)
		np := pl.GetPriority(coin.Price.Provider)
		if np > op {
			return errors.E("ticker provider with less priority")
		}

		if cd.LastUpdate.After(coin.LastUpdate) && op >= np {
			return errors.E("ticker is outdated")
		}
	}
	hm := createHashMap(coin.CoinName, coin.TokenId)
	return s.AddHM(EntityQuotes, hm, coin)
}

func (s *Storage) GetTicker(coin, token string) (blockatlas.Ticker, error) {
	hm := createHashMap(coin, token)
	var cd *blockatlas.Ticker
	err := s.GetHMValue(EntityQuotes, hm, &cd)
	if err != nil {
		return blockatlas.Ticker{}, err
	}
	return *cd, nil
}

type RateProviderList interface {
	GetPriority(providerId string) int
}

func (s *Storage) SaveRates(rates blockatlas.Rates, pl RateProviderList) {
	for _, rate := range rates {
		r, err := s.GetRate(rate.Currency)
		if err == nil {
			op := pl.GetPriority(r.Provider)
			np := pl.GetPriority(rate.Provider)
			if np > op {
				logger.Error("rate provider with less priority")
			}

			if rate.Timestamp < r.Timestamp && op >= np {
				logger.Error("rate is outdated")
			}
		}
		err = s.AddHM(EntityRates, rate.Currency, &rate)
		if err != nil {
			logger.Error(err, "SaveRates", logger.Params{"rate": rate})
		}
	}
}

func (s *Storage) GetRate(currency string) (rate *blockatlas.Rate, err error) {
	err = s.GetHMValue(EntityRates, currency, &rate)
	return
}

func createHashMap(coin, token string) string {
	if len(token) == 0 {
		return strings.ToUpper(coin)
	}
	return strings.ToUpper(strings.Join([]string{coin, token}, "_"))
}
