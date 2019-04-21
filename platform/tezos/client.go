package tezos

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

type Client struct {
	HTTPClient *http.Client
	RpcURL     string
}

func (c *Client) GetTxsOfAddress(address string) ([]Tx, error) {
	uri := fmt.Sprintf("%s/operations/%s?type=Transaction",
		c.RpcURL, url.PathEscape(address))
	httpRes, err := c.HTTPClient.Get(uri)
	if err != nil {
		logrus.WithError(err).Error("Tezos: Failed to get transactions")
		return nil, ErrSourceConn
	}

	if httpRes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status %s", httpRes.Status)
	}

	var res []Tx
	err = json.NewDecoder(httpRes.Body).Decode(&res)

	return res, nil
}
