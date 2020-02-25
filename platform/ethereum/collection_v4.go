package ethereum

import (
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"strings"
)

func (p *Platform) GetCollectionsV4(owner string) (blockatlas.CollectionPage, error) {
	collections, err := p.collectionsClient.GetCollections(owner)
	if err != nil {
		return nil, err
	}
	return NormalizeCollectionPageV4(collections, p.CoinIndex, owner), nil
}

func (p *Platform) GetCollectiblesV4(owner, collectibleID string) (blockatlas.CollectiblePage, error) {
	items, err := p.collectionsClient.GetCollectiblesV4(owner, collectibleID)
	if err != nil {
		return nil, err
	}
	return NormalizeCollectiblePageV4(items, p.CoinIndex), nil
}

func NormalizeCollectionPageV4(collections []Collection, coinIndex uint, owner string) (page blockatlas.CollectionPage) {
	for _, collection := range collections {
		item := NormalizeCollectionV4(collection, coinIndex, owner)
		page = append(page, item)
	}
	return page
}

func NormalizeCollectionV4(c Collection, coinIndex uint, owner string) blockatlas.Collection {
	return blockatlas.Collection{
		Name:            c.Name,
		Slug:            c.Slug,
		ImageUrl:        c.ImageUrl,
		Description:     c.Description,
		ExternalLink:    c.ExternalUrl,
		Total:           int(c.Total.Int64()),
		Id:              c.Slug,
		CategoryAddress: c.Slug,
		Address:         owner,
		Coin:            coinIndex,
	}
}

func NormalizeCollectiblePageV4(srcPage []Collectible, coinIndex uint) (page blockatlas.CollectiblePage) {
	for _, src := range srcPage {
		item := NormalizeCollectibleV4(src, coinIndex)
		if _, ok := supportedTypes[item.Type]; !ok {
			continue
		}
		page = append(page, item)
	}
	return page
}

func NormalizeCollectibleV4(a Collectible, coinIndex uint) blockatlas.Collectible {
	id := strings.Join([]string{a.AssetContract.Address, a.TokenId}, "-")
	return blockatlas.Collectible{
		ID:              id,
		CollectionID:    a.Collection.Slug,
		TokenID:         a.TokenId,
		ContractAddress: a.AssetContract.Address,
		Name:            a.Name,
		Category:        a.Collection.Name,
		ImageUrl:        a.ImagePreviewUrl,
		ProviderLink:    a.Permalink,
		ExternalLink:    a.ExternalLink,
		Type:            a.AssetContract.Type,
		Description:     a.Description,
		Coin:            coinIndex,
		Version:         a.AssetContract.Version,
	}
}
