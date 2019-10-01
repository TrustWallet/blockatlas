package icon

import (
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"net/url"
	"strconv"
)

type Client struct {
	blockatlas.Request
}

func InitClient(baseUrl string) Client {
	return Client{
		Request: blockatlas.Request{
			HttpClient:   blockatlas.DefaultClient,
			ErrorHandler: blockatlas.DefaultErrorHandler,
			BaseUrl:      baseUrl,
		},
	}
}

func (c *Client) GetAddressTransactions(address string) ([]Tx, error) {
	query := url.Values{
		"address": {address},
		"count":   {strconv.FormatInt(blockatlas.TxPerPage, 10)},
	}
	var res Response
	err := c.Get(&res, "/address/txList", query)
	if err != nil {
		return nil, err
	}
	return res.Data, nil
}
