package db

import (
	"context"
	"github.com/trustwallet/blockatlas/db/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

func (i Instance) GetSubscribedAddressesForAssets(ctx context.Context, addresses []string) ([]models.Address, error) {
	db := i.Gorm.WithContext(ctx)
	var result []models.Address
	err := db.Model(&models.AssetSubscription{}).
		Select("id", "address").
		Joins("LEFT JOIN addresses a ON a.id = address_id").
		Where("address in (?)", addresses).
		Scan(&result).
		Limit(len(addresses)).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i Instance) GetAssetsMapByAddresses(addresses []string, ctx context.Context) (map[string][]string, error) {
	db := i.Gorm.WithContext(ctx)
	var associations []models.AddressToAssetAssociation
	if err := db.Joins("Address").Joins("Asset").Find(&associations, "address in (?)", addresses).Error; err != nil {
		return nil, err
	}

	result := make(map[string][]string)
	for _, a := range associations {
		assets := result[a.Address.Address]
		result[a.Address.Address] = append(assets, a.Asset.Asset)
	}
	return result, nil
}

func (i Instance) GetAssetsMapByAddressesFromTime(addresses []string, from time.Time, ctx context.Context) (map[string][]string, error) {
	if len(addresses) == 0 {
		return map[string][]string{}, nil
	}
	db := i.Gorm.WithContext(ctx)
	var associations []models.AddressToAssetAssociation
	err := db.Joins("Address").Where("address in (?)", addresses).Joins("Asset").Find(&associations, "created_at > ?", from).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string][]string)
	for _, a := range associations {
		assets := result[a.Address.Address]
		result[a.Address.Address] = append(assets, a.Asset.Asset)
	}
	return result, nil
}

func (i *Instance) GetAssociationsByAddresses(addresses []string, ctx context.Context) ([]models.AddressToAssetAssociation, error) {
	db := i.Gorm.WithContext(ctx)
	var result []models.AddressToAssetAssociation
	if err := db.Joins("Address").Joins("Asset").Find(&result, "address in (?)", addresses).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (i *Instance) GetAssociationsByAddressesFromTime(addresses []string, from time.Time, ctx context.Context) ([]models.AddressToAssetAssociation, error) {
	db := i.Gorm.WithContext(ctx)
	var result []models.AddressToAssetAssociation
	err := db.Joins("Address").Where("address in (?)", addresses).Joins("Asset").Find(&result, "created_at > ?", from).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (i *Instance) AddAssociationsForAddress(address string, assets []string, ctx context.Context) error {
	db := i.Gorm.WithContext(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		uniqueAssets := getUniqueStrings(assets)
		uniqueAssetsModel := make([]models.Asset, 0, len(uniqueAssets))
		for _, l := range uniqueAssets {
			uniqueAssetsModel = append(uniqueAssetsModel, models.Asset{
				Asset: l,
			})
		}

		var err error
		dbAddress := models.Address{Address: address}
		err = db.Clauses(clause.OnConflict{DoNothing: true}).FirstOrCreate(&dbAddress, "address = ?", address).Error
		if err != nil {
			return err
		}

		if len(uniqueAssetsModel) > 0 {
			if err = db.Clauses(clause.OnConflict{DoNothing: true}).Create(&uniqueAssetsModel).Error; err != nil {
				return err
			}
		}

		var dbAssets []models.Asset
		if len(uniqueAssets) > 0 {
			if err = db.Where("asset in (?)", uniqueAssets).Find(&dbAssets).Error; err != nil {
				return err
			}
		}

		assetsSub := models.AssetSubscription{AddressID: dbAddress.ID}
		err = db.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{
					Name: "address_id",
				},
			},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"deleted_at": nil,
			}),
		}).Create(&assetsSub).Error
		if err != nil {
			return err
		}

		result := make([]models.AddressToAssetAssociation, 0, len(dbAssets))
		for _, asset := range dbAssets {
			result = append(result, models.AddressToAssetAssociation{
				AddressID: dbAddress.ID,
				AssetID:   asset.ID,
			})
		}
		if len(result) > 0 {
			return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&result).Error
		}
		return nil
	})
}

func (i *Instance) UpdateAssociationsForExistingAddresses(associations map[string][]string, ctx context.Context) error {
	db := i.Gorm.WithContext(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		assets := make([]string, 0, len(associations))
		for _, v := range associations {
			assets = append(assets, v...)
		}

		if len(assets) == 0 {
			return nil
		}

		uniqueAssets := getUniqueStrings(assets)
		uniqueAssetsModel := make([]models.Asset, 0, len(uniqueAssets))
		for _, l := range uniqueAssets {
			uniqueAssetsModel = append(uniqueAssetsModel, models.Asset{Asset: l})
		}
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&uniqueAssetsModel).Error; err != nil {
			return err
		}

		var dbAssets []models.Asset
		err := db.Where("asset in (?)", uniqueAssets).
			Find(&dbAssets).
			Limit(len(uniqueAssets)).
			Error
		if err != nil {
			return err
		}

		assetsMap := makeMapAssets(dbAssets)

		addresses := make([]string, 0, len(associations))
		for k := range associations {
			addresses = append(addresses, k)
		}

		var dbAddresses []models.Address
		if err := db.Find(&dbAddresses, "address in (?)", addresses).Limit(len(addresses)).Error; err != nil {
			return err
		}

		var addressSubs []models.AssetSubscription
		for _, a := range dbAddresses {
			sub := models.AssetSubscription{AddressID: a.ID}
			addressSubs = append(addressSubs, sub)
		}

		err = db.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{
					Name: "address_id",
				},
			},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"deleted_at": nil,
			}),
		}).Create(&addressSubs).Error
		if err != nil {
			return err
		}

		addressesMap := makeMapAddress(dbAddresses)

		var result []models.AddressToAssetAssociation
		for address, assets := range associations {
			for _, asset := range assets {
				addressID, ok := addressesMap[address]
				if !ok {
					continue
				}
				assetID, ok := assetsMap[asset]
				if !ok {
					continue
				}
				r := models.AddressToAssetAssociation{
					AddressID: addressID,
					AssetID:   assetID,
				}
				result = append(result, r)
			}
		}
		return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&result).Error
	})
}

func makeMapAssets(addresses []models.Asset) map[string]uint {
	result := make(map[string]uint)
	for _, a := range addresses {
		result[a.Asset] = a.ID
	}
	return result
}

func makeMapAddress(addresses []models.Address) map[string]uint {
	result := make(map[string]uint)
	for _, a := range addresses {
		result[a.Address] = a.ID
	}
	return result
}

func getUniqueStrings(values []string) []string {
	keys := make(map[string]struct{})
	var list []string
	for _, entry := range values {
		if _, value := keys[entry]; !value {
			keys[entry] = struct{}{}
			list = append(list, entry)
		}
	}
	return list
}
