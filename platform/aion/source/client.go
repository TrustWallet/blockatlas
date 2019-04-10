package source

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/trustwallet/blockatlas/models"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	HttpClient *http.Client
	RpcUrl     string
}

func (c *Client) GetTxsOfAddress(address string) ([]AionTransaction, error) {
	uri := fmt.Sprintf("%sgetTransactionsByAddress?%s",
		c.RpcUrl,
		url.Values{
			"accountAddress": {address},
			"size":           {strconv.FormatInt(models.TxPerPage, 10)},
		}.Encode())

	res, err := c.HttpClient.Get(uri)
	if err != nil {
		logrus.WithError(err).Errorf("Aion: Failed to get transactions for address %s", address)
	}
	defer res.Body.Close()

	var trxs AionTransactionRPCResponse
	err = json.NewDecoder(res.Body).Decode(&trxs)
	
	return trxs.Content, nil
}

