package stellar

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stellar/go/clients/horizon"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/models"
	"github.com/trustwallet/blockatlas/platform/stellar/source"
	"github.com/trustwallet/blockatlas/util"
	"net/http"
	"strconv"
	"time"
)

var stellarClient = source.Client{
	HTTP: &http.Client{
		Timeout: 2 * time.Second,
	},
}

func Setup(router gin.IRouter) {
	router.Use(util.RequireConfig("stellar.api"))
	router.Use(func(c *gin.Context) {
		stellarClient.API = viper.GetString("stellar.api")
		c.Next()
	})
	router.GET("/:address", func(c *gin.Context) {
		GetTransactions(c, &coin.XLM, &stellarClient)
	})
}

func GetTransactions(c *gin.Context, nativeCoin *coin.Coin, client *source.Client) {
	payments, err := client.GetTxsOfAddress(c.Param("address"))
	if apiError(c, err) {
		return
	}

	txs := make([]models.Tx, 0)
	for _, payment := range payments {
		tx, ok := FormatTx(&payment, nativeCoin)
		if !ok {
			continue
		}
		txs = append(txs, tx)
	}

	page := models.Response(txs)
	page.Sort()
	c.JSON(http.StatusOK, &page)
}

func apiError(c *gin.Context, err error) bool {
	if hErr, ok := err.(*horizon.Error); ok {
		if hErr.Problem.Type == "https://stellar.org/horizon-errors/bad_request" {
			c.String(http.StatusBadRequest, "Bad request!")
			return true
		} else {
			c.String(http.StatusBadRequest, hErr.Problem.Type)
			return true
		}
	}
	if err != nil {
		logrus.WithError(err).Warning("Stellar API request failed")
		c.String(http.StatusBadGateway, "Stellar API request failed")
		return true
	}
	return false
}

func FormatTx(payment *horizon.Payment, nativeCoin *coin.Coin) (tx models.Tx, ok bool) {
	if payment.Type != "payment" {
		return tx, false
	}
	if payment.AssetType != "native" {
		return tx, false
	}
	id, err := strconv.ParseUint(payment.ID, 10, 64)
	if err != nil {
		return tx, false
	}
	date, err := time.Parse("2006-01-02T15:04:05Z", payment.CreatedAt)
	if err != nil {
		return tx, false
	}
	value, err := util.DecimalToSatoshis(payment.Amount)
	if err != nil {
		return tx, false
	}
	value += "00" // 5 decimal places to 7
	return models.Tx{
		Id:    payment.TransactionHash,
		From:  payment.From,
		To:    payment.To,
		// https://www.stellar.org/developers/guides/concepts/fees.html
		// Fee fixed at 100 stroops
		Fee:   "100",
		Date:  date.Unix(),
		Block: id,
		Meta:  models.Transfer{
			Name:     nativeCoin.Title,
			Symbol:   nativeCoin.Symbol,
			Decimals: 7,
			Value:    value,
		},
	}, true
}
