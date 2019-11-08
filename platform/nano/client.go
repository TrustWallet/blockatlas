package nano

import (
	"strconv"

	"github.com/trustwallet/blockatlas/pkg/blockatlas"
)

type Client struct {
	blockatlas.Request
}

func (c *Client) GetAccountHistory(address string) (history AccountHistory, err error) {
	count := strconv.Itoa(blockatlas.TxPerPage)
	err = c.Post(&history, "/", AccountHistoryRequest{Action: "account_history", Account: address, Count: count})
	return history, err
}
