package cmc

import (
	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas/marketdata/cmcmap"
	"github.com/trustwallet/blockatlas/marketdata/market"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"math/big"
	"net/url"
)

type Market struct {
	market.Market
}

func InitMarket() market.Provider {
	m := &Market{
		Market: market.Market{
			Id:         "cmc",
			Request:    blockatlas.InitClient(viper.GetString("market.cmc_api")),
			UpdateTime: viper.GetString("market.cmc_quote_update_time"),
		},
	}
	m.Headers["X-CMC_PRO_API_KEY"] = viper.GetString("market.cmc_api_key")
	return m
}

func (m *Market) GetData() (blockatlas.Tickers, error) {
	cmap, err := cmcmap.GetCmcMap()
	if err != nil {
		return nil, err
	}
	var prices CoinPrices
	err = m.Get(&prices, "v1/cryptocurrency/listings/latest", url.Values{"limit": {"5000"}, "convert": {"USD"}})
	if err != nil {
		return nil, err
	}
	return normalizeTickers(prices, m.GetId(), cmap), nil
}

func normalizeTicker(price Data, provider string, cmap cmcmap.CmcMapping) (blockatlas.Ticker, error) {
	tokenId := ""
	coinName := price.Symbol
	coinType := blockatlas.TypeCoin
	if price.Platform != nil {
		coinType = blockatlas.TypeToken
		coinName = price.Platform.Symbol
		tokenId = price.Platform.TokenAddress
		if len(tokenId) == 0 {
			tokenId = price.Symbol
		}
	}
	cmcCoin, cmcTokenId, err := cmap.GetCoin(price.Id)
	if err == nil {
		coinName = cmcCoin.Symbol
		tokenId = cmcTokenId
	}
	return blockatlas.Ticker{
		CoinName: coinName,
		CoinType: coinType,
		TokenId:  tokenId,
		Price: blockatlas.TickerPrice{
			Value:     big.NewFloat(price.Quote.USD.Price),
			Change24h: big.NewFloat(price.Quote.USD.PercentChange24h),
			Currency:  "USD",
			Provider:  provider,
		},
		LastUpdate: price.LastUpdated,
	}, nil
}

func normalizeTickers(prices CoinPrices, provider string, cmap cmcmap.CmcMapping) (tickers blockatlas.Tickers) {
	for _, price := range prices.Data {
		t, err := normalizeTicker(price, provider, cmap)
		if err != nil {
			logger.Error(err)
			continue
		}
		tickers = append(tickers, t)
	}
	return
}

func percentageChange(value, percent float64) float64 {
	return value * (percent / 100)
}
