package api

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/ginutils"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/blockatlas/platform"
	"strconv"
	"sync"
)

type AddressBatchRequest struct {
	Coin    uint   `json:"coin"`
	Address string `json:"address"`
}

type ENSBatchRequest struct {
	Coins []uint64 `json:"coins"`
	Name  string   `json:"name"`
}

type AddressesRequest []AddressBatchRequest

// @Summary Get Multiple Stake Delegations
// @ID batch_delegations
// @Description Get Stake Delegations for multiple coins
// @Accept json
// @Produce json
// @Tags platform,staking
// @Param delegations body api.AddressesRequest true "Validators addresses and coins"
// @Success 200 {object} blockatlas.DelegationsBatchPage
// @Router /v2/staking/delegations [post]
func makeStakingDelegationsBatchRoute(router gin.IRouter) {
	router.POST("/staking/delegations", func(c *gin.Context) {
		var reqs AddressesRequest
		if err := c.BindJSON(&reqs); err != nil {
			ginutils.ErrorResponse(c).Message(err.Error()).Render()
			return
		}

		batch := make(blockatlas.DelegationsBatchPage, 0)
		for _, r := range reqs {
			c := coin.Coins[r.Coin]
			p := platform.StakeAPIs[c.Handle]
			batch = append(batch, getDelegationResponse(p, r.Address))
		}
		ginutils.RenderSuccess(c, blockatlas.DocsResponse{Docs: batch})
	})
}

// @Description Get collection categories
// @ID collection_categories
// @Summary Get list of collections from a specific coin and addresses
// @Accept json
// @Produce json
// @Tags Collectibles
// @Param data body string true "Payload" default({"60": ["0xb3624367b1ab37daef42e1a3a2ced012359659b0"]})
// @Success 200 {object} blockatlas.DocsResponse
// @Router /v2/collectibles/categories [post]
func makeCategoriesBatchRoute(router gin.IRouter) {
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

// @Description Get multiple addresses from naming service
// @ID naming_service
// @Summary Get list of addresses for passed coins
// @Accept json
// @Produce json
// @Tags ns
// @Param coins body api.ENSBatchRequest true "ns name and coins"
// @Success 200 {object} api.LookupBatchPage
// @Router /v2/ns/lookup [post]
func makeNsLookupRoute(router gin.IRouter) {
	router.POST("/ns/lookup", func(c *gin.Context) {
		var req ENSBatchRequest
		if err := c.BindJSON(&req); err != nil {
			ginutils.ErrorResponse(c).Message(err.Error()).Render()
			return
		}

		lookupsChan := make(chan blockatlas.Resolved)
		var wg sync.WaitGroup
		wg.Add(len(req.Coins))

		for _, coin := range req.Coins {
			go func(coin uint64) {
				defer wg.Done()
				lookup, err := handleLookup(req.Name, coin)
				if err != nil {
					logger.Error(err, logger.Params{"name": req.Name, "coin": coin, })
					return
				}
				lookupsChan <- lookup
			}(coin)
		}

		// extract data from the channel
		batch := make(LookupBatchPage, 0)
		go func() {
			for lookup := range lookupsChan {
				batch = append(batch, lookup)
			}
		}()

		wg.Wait()
		close(lookupsChan)

		ginutils.RenderSuccess(c, batch)
	})
}
