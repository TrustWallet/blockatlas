package nano

import (
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"

	"github.com/spf13/viper"
)

type Platform struct {
	client Client
}

func (p *Platform) Init() error {
	p.client = Client{blockatlas.InitClient(viper.GetString("nano.api"))}
	p.client.Headers["Content-Type"] = "application/json"
	return nil
}

func (p *Platform) Coin() coin.Coin {
	return coin.Coins[coin.NANO]
}
