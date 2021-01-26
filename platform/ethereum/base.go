package ethereum

import (
	"github.com/trustwallet/blockatlas/platform/bitcoin/blockbook"
	"github.com/trustwallet/blockatlas/platform/ethereum/bounce"
	"github.com/trustwallet/blockatlas/platform/ethereum/opensea"
	"github.com/trustwallet/blockatlas/platform/ethereum/trustray"
	"github.com/trustwallet/golibs/client"
	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/golibs/network/middleware"
)

type Platform struct {
	CoinIndex   uint
	client      EthereumClient
	collectible CollectibleClient
}

func InitTrustRay(coinType uint, api string) *Platform {
	return &Platform{
		CoinIndex: coinType,
		client:    &trustray.Client{Request: client.InitClient(api, middleware.SentryErrorHandler)},
	}
}

func InitWithBlockbook(coinType uint, blockbookApi string) *Platform {
	return &Platform{
		CoinIndex: coinType,
		client:    &blockbook.Client{Request: client.InitClient(blockbookApi, middleware.SentryErrorHandler)},
	}
}

func InitWithOpenSea(coinType uint, blockbookApi, collectionApi, collectionKey string) *Platform {
	platform := InitWithBlockbook(coinType, blockbookApi)
	platform.collectible = opensea.InitClient(collectionApi, collectionKey)
	return platform
}

func InitWithBounce(coinType uint, blockbookApi, collectionApi string) *Platform {
	platform := InitWithBlockbook(coinType, blockbookApi)
	platform.collectible = bounce.InitClient(collectionApi)
	return platform
}

func (p *Platform) Coin() coin.Coin {
	return coin.Coins[p.CoinIndex]
}
