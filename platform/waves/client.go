package waves

import (
	"fmt"
	"github.com/trustwallet/blockatlas/client"
	"net/http"
	"net/url"
)

type Client struct {
	HTTPClient *http.Client
	URL        string
}

func (c *Client) GetTxs(address string, limit int) ([]Transaction, error) {
	path := fmt.Sprintf("transactions/address/%s/limit/%d", address, limit)

	txsArrays := make([][]Transaction, 0)
	err := client.Request(c.HTTPClient, c.URL, path, url.Values{}, &txsArrays)

	if len(txsArrays) > 0 {
		return txsArrays[0], err
	} else {
		return []Transaction{}, err
	}
}

func (c *Client) GetBlockByNumber(num int64) (*Block, error) {
	path := fmt.Sprintf("blocks/at/%d", num)

	block := new(Block)
	err := client.Request(c.HTTPClient, c.URL, path, url.Values{}, &block)

	return block, err
}

func (c *Client) GetCurrentBlock() (*CurrentBlock, error) {
	path := "blocks/height"

	block := new(CurrentBlock)
	err := client.Request(c.HTTPClient, c.URL, path, url.Values{}, &block)

	return block, err
}
