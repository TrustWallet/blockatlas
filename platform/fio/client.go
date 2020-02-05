package fio

import (
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/errors"
)

// Client for FIO API
type Client struct {
	blockatlas.Request
}

func (c *Client) lookupPubAddress(name string, coinSymbol string) (address string, error error) {
	var res GetPubAddressResponse
	err := c.Post(&res, "get_pub_address", GetPubAddressRequest{FioAddress: name, TokenCode: coinSymbol})
	if err != nil {
		return "", errors.E(err, "Error lokking up FIO name", errors.Params{"name": name, "coinSymbol": coinSymbol, "inner_error": err.Error()})
	}
	if res.Message != "" {
		return "", errors.E("Error lokking up FIO name", errors.Params{"name": name, "coinSymbol": coinSymbol, "inner_error": res.Message})
	}
	return res.PublicAddress, nil
}
