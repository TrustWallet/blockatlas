package api

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/ginutils"
	"github.com/trustwallet/blockatlas/platform"
	"strconv"
)

// @Description Get collection categories
// @ID collection_categories_v4
// @Summary Get list of collections from a specific coin and addresses
// @Accept json
// @Produce json
// @Tags Collections
// @Param data body string true "Payload" default({"60": ["0xb3624367b1ab37daef42e1a3a2ced012359659b0"]})
// @Success 200 {object} blockatlas.DocsResponse
// @Router /v4/collectibles/categories [post]
func makeCategoriesBatchRouteV4(router gin.IRouter) {
	router.POST("/collectibles/categories", func(c *gin.Context) {
		var reqs map[string][]string
		if err := c.BindJSON(&reqs); err != nil {
			ginutils.ErrorResponse(c).Message(err.Error()).Render()
			return
		}

		batch := make(blockatlas.CollectionPage, 0)
		for key, addresses := range reqs {
			coinId, err := strconv.Atoi(key)
			if err != nil {
				continue
			}
			p, ok := platform.CollectionAPIs[uint(coinId)]
			if !ok {
				continue
			}
			for _, address := range addresses {
				collections, err := p.GetCollections(address)
				if err != nil {
					continue
				}
				batch = append(batch, collections...)
			}
		}
		ginutils.RenderSuccess(c, batch)
	})
}

// @Summary Get Collection
// @ID collection_v4
// @Description Get a collection from the address
// @Accept json
// @Produce json
// @Tags Collections
// @Param coin path string true "the coin name" default(ethereum)
// @Param owner path string true "the query address" default(0x0875BCab22dE3d02402bc38aEe4104e1239374a7)
// @Param collection_id path string true "the query collection" default(0x06012c8cf97bead5deae237070f9587f8e7a266d)
// @Success 200 {object} blockatlas.CollectionPage
// @Failure 500 {object} ginutils.ApiError
// @Router /v4/{coin}/collections/{owner}/collection/{collection_id} [get]
func makeCollectionRouteV4(router gin.IRouter, api blockatlas.Platform) {
	var collectionAPI blockatlas.CollectionAPI
	collectionAPI, _ = api.(blockatlas.CollectionAPI)

	if collectionAPI == nil {
		return
	}

	router.GET("/collections/:owner/collection/:collection_id", func(c *gin.Context) {
		collectibles, err := collectionAPI.GetCollectibles(c.Param("owner"), c.Param("collection_id"))
		if err != nil {
			ginutils.ErrorResponse(c).Message(err.Error()).Render()
			return
		}

		ginutils.RenderSuccess(c, collectibles)
	})
}

func emptyPage(c *gin.Context) {
	var page blockatlas.TxPage
	ginutils.RenderSuccess(c, &page)
}
