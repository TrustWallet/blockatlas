package aion

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/models"
	"github.com/trustwallet/blockatlas/platform/aion/source"
	"github.com/trustwallet/blockatlas/util"
	"net/http"
	"strconv"
)

var client = source.Client{
	HttpClient: http.DefaultClient,
}

func Setup(router gin.IRouter) {
	router.Use(util.RequireConfig("aion.api"))
	router.Use(func(c *gin.Context) {
		client.RpcUrl = viper.GetString("aion.api")
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
		fee := strconv.Itoa(srcTx.NrgConsumed)

		txs = append(txs, models.Tx{
			Id:    srcTx.BlockHash,
			Coin:  coin.IndexAION,
			Date:  srcTx.BlockTimestamp,
			From:  "0x" + srcTx.FromAddr,
			To:    "0x" + srcTx.ToAddr,
			Fee:   models.Amount(fee),
			Block: srcTx.BlockNumber,
			Meta:  models.Transfer{
				Value: srcTx.Value,
			},
		})
	}

	page := models.Response(txs)
	page.Sort()
	c.JSON(http.StatusOK, &page)
}
