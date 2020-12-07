package tokensearcher

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/blockatlas/db"
	"github.com/trustwallet/blockatlas/db/models"
	"github.com/trustwallet/blockatlas/mq"
	"github.com/trustwallet/blockatlas/pkg/address"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"strconv"
	"sync"
	"time"
)

type (
	AddressesByCoin map[uint][]string
	AssetsByAddress map[string][]models.Asset
	Request         struct {
		AddressesByCoin map[string][]string `json:"addresses"`
		From            uint                `json:"from"`
	}
)

type Instance struct {
	database *db.Instance
	apis     map[uint]blockatlas.TokensAPI
	mqClient *new_mq.Client
	queue    new_mq.Queue
}

func Init(database *db.Instance, apis map[uint]blockatlas.TokensAPI, mqClient *new_mq.Client, queue new_mq.Queue) Instance {
	return Instance{database: database, apis: apis, mqClient: mqClient, queue: queue}
}

func (i Instance) HandleTokensRequest(request Request, ctx context.Context) (map[string][]string, error) {
	addresses := getAddressesFromRequest(request)
	if len(addresses) == 0 {
		return nil, nil
	}

	subscribedAddresses, err := getSubscribedAddresses(i.database, addresses, ctx)
	if err != nil {
		return nil, err
	}
	log.Info("subscribedAddresses " + strconv.Itoa(len(subscribedAddresses)))
	unsubscribedAddresses := getUnsubscribedAddresses(subscribedAddresses, addresses)

	log.Info("unsubscribedAddresses " + strconv.Itoa(len(unsubscribedAddresses)))
	assetsFromDB, err := i.database.GetAssetsMapByAddressesFromTime(
		subscribedAddresses,
		time.Unix(int64(request.From), 0),
		ctx)
	if err != nil {
		return nil, err
	}

	log.Info("assetsFromDB " + strconv.Itoa(len(assetsFromDB)))
	assetsFromNodes := make(AssetsByAddress)
	if len(unsubscribedAddresses) != 0 {
		assetsFromNodes = getAssetsByAddressFromNodes(unsubscribedAddresses, i.apis)
		err = i.publishNewAddressesToQueue(assetsFromNodes)
		if err != nil {
			log.Error(err)
		}
	}

	return getAssetsToResponse(assetsFromDB, assetsFromNodes, addresses), nil
}

func getSubscribedAddresses(database *db.Instance, addresses []string, ctx context.Context) ([]string, error) {
	subscribedAddressesModel, err := database.GetSubscribedAddressesForAssets(ctx, addresses)
	if err != nil {
		return nil, err
	}

	subscribedAddresses := make([]string, 0, len(subscribedAddressesModel))
	for _, a := range subscribedAddressesModel {
		subscribedAddresses = append(subscribedAddresses, a.Address)
	}
	return subscribedAddresses, nil
}

func getAddressesFromRequest(request Request) []string {
	var addresses []string
	for coinID, requestAddresses := range request.AddressesByCoin {
		for _, a := range requestAddresses {
			addresses = append(addresses, coinID+"_"+a)
		}
	}
	return addresses
}

func getUnsubscribedAddresses(subscribed []string, all []string) AddressesByCoin {
	addressesByCoin := make(AddressesByCoin)
	subscribedMap := make(map[string]bool)
	for _, a := range subscribed {
		subscribedMap[a] = true
	}
	for _, a := range all {
		_, ok := subscribedMap[a]
		if !ok {
			ua, coinID, ok := address.UnprefixedAddress(a)
			if !ok {
				continue
			}
			currentAddresses := addressesByCoin[coinID]
			addressesByCoin[coinID] = append(currentAddresses, ua)
		}
	}
	return addressesByCoin
}

func getAddressesToRegisterByCoin(assetsByAddresses AssetsByAddress, addresses []string) AddressesByCoin {
	addressesByCoin := make(AddressesByCoin)
	addressesFromRequestMap := make(map[string]bool)
	for _, a := range addresses {
		addressesFromRequestMap[a] = true
	}
	for _, a := range addresses {
		_, ok := assetsByAddresses[a]
		if !ok {
			ua, coinID, ok := address.UnprefixedAddress(a)
			if !ok {
				continue
			}
			currentAddresses := addressesByCoin[coinID]
			addressesByCoin[coinID] = append(currentAddresses, ua)
		}
	}
	return addressesByCoin
}

func getAssetsByAddressFromNodes(addressesByCoin AddressesByCoin, apis map[uint]blockatlas.TokensAPI) AssetsByAddress {
	a := NodesResponse{AssetsByAddress: make(AssetsByAddress)}
	var wg sync.WaitGroup
	for coinID, addresses := range addressesByCoin {
		api, ok := apis[coinID]
		if !ok {
			continue
		}
		wg.Add(1)
		go fetchAssetsByAddresses(api, addresses, &a, &wg)
	}
	wg.Wait()
	return a.AssetsByAddress
}

func fetchAssetsByAddresses(tokenAPI blockatlas.TokensAPI, addresses []string, result *NodesResponse, wg *sync.WaitGroup) {
	defer wg.Done()

	var tWg sync.WaitGroup
	tWg.Add(len(addresses))
	for _, a := range addresses {
		go func(address string, tWg *sync.WaitGroup) {
			defer tWg.Done()
			tokens, err := tokenAPI.GetTokenListByAddress(address)
			if err != nil {
				log.Error("Chain: " + tokenAPI.Coin().Handle + " Address: " + address)
				return
			}
			result.UpdateAssetsByAddress(tokens, int(tokenAPI.Coin().ID), address)
		}(a, &tWg)
	}
	tWg.Wait()
}

func (i *Instance) publishNewAddressesToQueue(message AssetsByAddress) error {
	log.Info("Published to queue")
	body, err := json.Marshal(message)
	log.Info(string(body))
	if err != nil {
		return err
	}
	return i.mqClient.Push(i.queue, body)
}

func getAssetsToResponse(dbAssetsMap, nodesAssetsMap AssetsByAddress, assetAddresses []string) map[string][]string {
	result := make(map[string][]string)
	for _, assetAddress := range assetAddresses {
		dbAddresses, ok := dbAssetsMap[assetAddress]
		if !ok {
			nodesAssets, ok := nodesAssetsMap[assetAddress]
			if !ok {
				continue
			}
			result[assetAddress] = models.AssetIDs(nodesAssets)
			continue
		}
		result[assetAddress] = models.AssetIDs(dbAddresses)
	}
	return result
}
