package aion

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/models"
	"github.com/trustwallet/blockatlas/util"
	"net/http"
	"strconv"
)

var client = Client{
	HTTPClient: http.DefaultClient,
}

// Setup registers the Aion route
func Setup(router gin.IRouter) {
	router.Use(util.RequireConfig("aion.api"))
	router.Use(func(c *gin.Context) {
		client.BaseURL = viper.GetString("aion.api")
		c.Next()
	})
	router.GET("/:address", getTransactions)
}

func getTransactions(c *gin.Context) {
	srcPage, err := client.GetTxsOfAddress(c.Param("address"))

	if err != nil {
		logrus.WithError(err).
			Errorf("Aion: Failed to get transactions for %s",
				c.Param("address"))
	}

	var txs []models.Tx
	for _, srcTx := range srcPage.Content {
		txs = append(txs, Normalize(&srcTx))
	}

	page := models.Response(txs)
	page.Sort()
	c.JSON(http.StatusOK, &page)
}

// Normalize converts an Aion transaction into the generic model
func Normalize(srcTx *Tx) models.Tx {
	fee := strconv.Itoa(srcTx.NrgConsumed)
	value := util.DecimalExp(string(srcTx.Value), 18)

	return models.Tx{
		ID:    srcTx.TransactionHash,
		Coin:  coin.AION,
		Date:  srcTx.BlockTimestamp,
		From:  "0x" + srcTx.FromAddr,
		To:    "0x" + srcTx.ToAddr,
		Fee:   models.Amount(fee),
		Block: srcTx.BlockNumber,
		Meta:  models.Transfer{
			Value: models.Amount(value),
		},
	}
}
