// +build integration

package db_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/blockatlas/db"
	"github.com/trustwallet/blockatlas/db/models"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/services/observer/subscriber"
	"github.com/trustwallet/blockatlas/tests/integration/setup"
	"testing"
)

func TestDb_AddSubscriptions(t *testing.T) {
	setup.CleanupPgContainer(dbConn)
	var subscriptions []models.SubscriptionData
	id := uint(1)
	subscriptions = append(subscriptions, models.SubscriptionData{
		SubscriptionId: id,
		Coin:           60,
		Address:        "testAddr",
	})

	subscriptions = append(subscriptions, models.SubscriptionData{
		SubscriptionId: id,
		Coin:           61,
		Address:        "testAddr2",
	})

	subscriptions = append(subscriptions, models.SubscriptionData{
		SubscriptionId: id,
		Coin:           62,
		Address:        "testAddr3",
	})

	assert.Nil(t, db.AddSubscriptions(dbConn, id, subscriptions))

	subs, err := db.GetSubscriptionData(dbConn, 60, []string{"testAddr"})
	assert.Nil(t, err)
	assert.NotNil(t, subs)
	assert.Equal(t, 1, len(subs))
	assert.Equal(t, subscriptions[0].SubscriptionId, subs[0].SubscriptionId)
	assert.Equal(t, subscriptions[0].Coin, subs[0].Coin)
	assert.Equal(t, subscriptions[0].Address, subs[0].Address)

	subs, err = db.GetSubscriptionData(dbConn, 61, []string{"testAddr2"})
	assert.Nil(t, err)
	assert.NotNil(t, subs)
	assert.Equal(t, 1, len(subs))
	assert.Equal(t, subscriptions[1].SubscriptionId, subs[0].SubscriptionId)
	assert.Equal(t, subscriptions[1].Coin, subs[0].Coin)
	assert.Equal(t, subscriptions[1].Address, subs[0].Address)

	subs, err = db.GetSubscriptionData(dbConn, 62, []string{"testAddr3"})
	assert.Nil(t, err)
	assert.NotNil(t, subs)
	assert.Equal(t, 1, len(subs))
	assert.Equal(t, subscriptions[2].SubscriptionId, subs[0].SubscriptionId)
	assert.Equal(t, subscriptions[2].Coin, subs[0].Coin)
	assert.Equal(t, subscriptions[2].Address, subs[0].Address)
}

func TestDb_AddSubscriptionsWithRewrite(t *testing.T) {
	setup.CleanupPgContainer(dbConn)

	id := uint(1)

	var subscriptions []models.SubscriptionData
	subscriptions = append(subscriptions, models.SubscriptionData{
		SubscriptionId: id,
		Coin:           60,
		Address:        "testAddr",
	})

	subscriptions = append(subscriptions, models.SubscriptionData{
		SubscriptionId: id,
		Coin:           714,
		Address:        "testAddr",
	})

	subscriptions = append(subscriptions, models.SubscriptionData{
		SubscriptionId: id,
		Coin:           144,
		Address:        "testAddr",
	})

	assert.Nil(t, db.AddSubscriptions(dbConn, id, subscriptions))

	subs60, err := db.GetSubscriptionData(dbConn, 60, []string{"testAddr"})
	assert.Nil(t, err)
	assert.NotNil(t, subs60)
	assert.Equal(t, 1, len(subs60))
	assert.Equal(t, subscriptions[0].SubscriptionId, subs60[0].SubscriptionId)
	assert.Equal(t, subscriptions[0].Coin, subs60[0].Coin)
	assert.Equal(t, subscriptions[0].Address, subs60[0].Address)

	subs714, err := db.GetSubscriptionData(dbConn, 714, []string{"testAddr"})
	assert.Nil(t, err)
	assert.NotNil(t, subs714)
	assert.Equal(t, 1, len(subs714))
	assert.Equal(t, subscriptions[1].SubscriptionId, subs714[0].SubscriptionId)
	assert.Equal(t, subscriptions[1].Coin, subs714[0].Coin)
	assert.Equal(t, subscriptions[1].Address, subs714[0].Address)

	subs144, err := db.GetSubscriptionData(dbConn, 144, []string{"testAddr"})
	assert.Nil(t, err)
	assert.NotNil(t, subs144)
	assert.Equal(t, 1, len(subs144))
	assert.Equal(t, subscriptions[2].SubscriptionId, subs144[0].SubscriptionId)
	assert.Equal(t, subscriptions[2].Coin, subs144[0].Coin)
	assert.Equal(t, subscriptions[2].Address, subs144[0].Address)

	subscriptions = append(subscriptions, models.SubscriptionData{
		SubscriptionId: id,
		Coin:           60,
		Address:        "testAddr2",
	})

	subscriptions = append(subscriptions, models.SubscriptionData{
		SubscriptionId: id,
		Coin:           714,
		Address:        "testAddr2",
	})

	subscriptions = append(subscriptions, models.SubscriptionData{
		SubscriptionId: id,
		Coin:           144,
		Address:        "testAddr2",
	})

	assert.Nil(t, db.AddSubscriptions(dbConn, id, subscriptions))

	subs2N60, err := db.GetSubscriptionData(dbConn, 60, []string{"testAddr2"})
	assert.Nil(t, err)
	assert.Nil(t, err)
	assert.NotNil(t, subs2N60)
	assert.Equal(t, 1, len(subs2N60))
	assert.Equal(t, subscriptions[3].SubscriptionId, subs2N60[0].SubscriptionId)
	assert.Equal(t, subscriptions[3].Coin, subs2N60[0].Coin)
	assert.Equal(t, subscriptions[3].Address, subs2N60[0].Address)

	subs2N714, err := db.GetSubscriptionData(dbConn, 714, []string{"testAddr2"})
	assert.Nil(t, err)
	assert.Nil(t, err)
	assert.NotNil(t, subs2N714)
	assert.Equal(t, 1, len(subs2N714))
	assert.Equal(t, subscriptions[4].SubscriptionId, subs2N714[0].SubscriptionId)
	assert.Equal(t, subscriptions[4].Coin, subs2N714[0].Coin)
	assert.Equal(t, subscriptions[4].Address, subs2N714[0].Address)

	subs2N114, err := db.GetSubscriptionData(dbConn, 144, []string{"testAddr2"})
	assert.Nil(t, err)
	assert.Nil(t, err)
	assert.NotNil(t, subs2N114)
	assert.Equal(t, 1, len(subs2N114))
	assert.Equal(t, subscriptions[5].SubscriptionId, subs2N114[0].SubscriptionId)
	assert.Equal(t, subscriptions[5].Coin, subs2N114[0].Coin)
	assert.Equal(t, subscriptions[5].Address, subs2N114[0].Address)
}

func TestDb_FindSubscriptions(t *testing.T) {
	setup.CleanupPgContainer(dbConn)

	var subscriptionsA []blockatlas.Subscription
	id := uint(1)
	subscriptionsA = append(subscriptionsA, blockatlas.Subscription{
		Id:      id,
		Coin:    60,
		Address: "etherAddress",
	})

	subscriptionsA = append(subscriptionsA, blockatlas.Subscription{
		Id:      id,
		Coin:    714,
		Address: "binanceAddress",
	})

	subscriptionsA = append(subscriptionsA, blockatlas.Subscription{
		Id:      id,
		Coin:    148,
		Address: "AtomAddress",
	})

	subscriptionsA = append(subscriptionsA, blockatlas.Subscription{
		Id:      id,
		Coin:    144,
		Address: "XLMAddress",
	})

	subscriptionsA = append(subscriptionsA, blockatlas.Subscription{
		Id:      id,
		Coin:    61,
		Address: "ETCAddress",
	})

	assert.Nil(t, db.AddSubscriptions(dbConn, id, subscriber.ToSubscriptionData(subscriptionsA)))

	var subscriptionsB []blockatlas.Subscription

	for _, sub := range subscriptionsA {
		sub.Id = uint(2)
		subscriptionsB = append(subscriptionsB, sub)
	}
	assert.Nil(t, db.AddSubscriptions(dbConn, 2, subscriber.ToSubscriptionData(subscriptionsB)))

	returnedSubs, err := db.GetSubscriptionData(dbConn, 60, []string{"etherAddress"})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(returnedSubs))

	returnedSubs, err = db.GetSubscriptionData(dbConn, 714, []string{"binanceAddress"})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(returnedSubs))

	returnedSubs, err = db.GetSubscriptionData(dbConn, 144, []string{"XLMAddress"})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(returnedSubs))

	returnedSubs, err = db.GetSubscriptionData(dbConn, 148, []string{"AtomAddress"})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(returnedSubs))

	returnedSubs, err = db.GetSubscriptionData(dbConn, 61, []string{"ETCAddress"})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(returnedSubs))
}

func TestDb_DeleteSubscriptions(t *testing.T) {
	setup.CleanupPgContainer(dbConn)

	var subscriptions []models.SubscriptionData

	id := uint(1)
	subscriptions = append(subscriptions, models.SubscriptionData{
		SubscriptionId: id,
		Coin:           60,
		Address:        "testAddr",
	})

	subscriptions = append(subscriptions, models.SubscriptionData{
		SubscriptionId: id,
		Coin:           714,
		Address:        "testAddr2",
	})

	subscriptions = append(subscriptions, models.SubscriptionData{
		SubscriptionId: id,
		Coin:           144,
		Address:        "testAddr3",
	})

	assert.Nil(t, db.AddSubscriptions(dbConn, id, subscriptions))

	subs60, err := db.GetSubscriptionData(dbConn, 60, []string{"testAddr"})
	assert.Nil(t, err)
	assert.NotNil(t, subs60)
	assert.Equal(t, 1, len(subs60))
	assert.Equal(t, subscriptions[0].SubscriptionId, subs60[0].SubscriptionId)
	assert.Equal(t, subscriptions[0].Coin, subs60[0].Coin)
	assert.Equal(t, subscriptions[0].Address, subs60[0].Address)

	subs714, err := db.GetSubscriptionData(dbConn, 714, []string{"testAddr2"})
	assert.Nil(t, err)
	assert.NotNil(t, subs714)
	assert.Equal(t, 1, len(subs714))
	assert.Equal(t, subscriptions[1].SubscriptionId, subs714[0].SubscriptionId)
	assert.Equal(t, subscriptions[1].Coin, subs714[0].Coin)
	assert.Equal(t, subscriptions[1].Address, subs714[0].Address)

	subs144, err := db.GetSubscriptionData(dbConn, 144, []string{"testAddr3"})
	assert.Nil(t, err)
	assert.NotNil(t, subs144)
	assert.Equal(t, 1, len(subs144))
	assert.Equal(t, subscriptions[2].SubscriptionId, subs144[0].SubscriptionId)
	assert.Equal(t, subscriptions[2].Coin, subs144[0].Coin)
	assert.Equal(t, subscriptions[2].Address, subs144[0].Address)

	subsToDel := []models.SubscriptionData{subscriptions[0]}

	assert.Nil(t, db.DeleteSubscriptions(dbConn, subsToDel))

	subs714N2, err := db.GetSubscriptionData(dbConn, 714, []string{"testAddr2"})
	assert.Nil(t, err)
	assert.NotNil(t, subs714N2)
	assert.Equal(t, 1, len(subs714N2))
	assert.Equal(t, subscriptions[1].SubscriptionId, subs714N2[0].SubscriptionId)
	assert.Equal(t, subscriptions[1].Coin, subs714N2[0].Coin)
	assert.Equal(t, subscriptions[1].Address, subs714N2[0].Address)

	subs144N2, err := db.GetSubscriptionData(dbConn, 144, []string{"testAddr3"})
	assert.Nil(t, err)
	assert.NotNil(t, subs144N2)
	assert.Equal(t, 1, len(subs144N2))
	assert.Equal(t, subscriptions[2].SubscriptionId, subs144N2[0].SubscriptionId)
	assert.Equal(t, subscriptions[2].Coin, subs144N2[0].Coin)
	assert.Equal(t, subscriptions[2].Address, subs144N2[0].Address)

	subs60N2, err := db.GetSubscriptionData(dbConn, 60, []string{"testAddr"})
	assert.Nil(t, err)
	assert.Len(t, subs60N2, 0)
}

func TestDeleteAll(t *testing.T) {
	setup.CleanupPgContainer(dbConn)

	var subscriptions []models.SubscriptionData

	id := uint(1)
	subscriptions = append(subscriptions, models.SubscriptionData{
		SubscriptionId: id,
		Coin:           60,
		Address:        "testAddr",
	})

	subscriptions = append(subscriptions, models.SubscriptionData{
		SubscriptionId: id,
		Coin:           714,
		Address:        "testAddr2",
	})

	subscriptions = append(subscriptions, models.SubscriptionData{
		SubscriptionId: id,
		Coin:           144,
		Address:        "testAddr3",
	})
	assert.Nil(t, db.AddSubscriptions(dbConn, id, subscriptions))

	subs60, err := db.GetSubscriptionData(dbConn, 60, []string{"testAddr"})
	assert.Nil(t, err)
	assert.Len(t, subs60, 1)

	subs714, err := db.GetSubscriptionData(dbConn, 714, []string{"testAddr2"})
	assert.Nil(t, err)
	assert.Len(t, subs714, 1)

	subs144, err := db.GetSubscriptionData(dbConn, 144, []string{"testAddr3"})
	assert.Nil(t, err)
	assert.Len(t, subs144, 1)

	assert.Nil(t, db.DeleteAllSubscriptions(dbConn, 1))

	subs60, err = db.GetSubscriptionData(dbConn, 60, []string{"testAddr"})
	assert.Nil(t, err)
	assert.Len(t, subs60, 0)

	subs714, err = db.GetSubscriptionData(dbConn, 714, []string{"testAddr2"})
	assert.Nil(t, err)
	assert.Len(t, subs714, 0)

	subs144, err = db.GetSubscriptionData(dbConn, 144, []string{"testAddr3"})
	assert.Nil(t, err)
	assert.Len(t, subs144, 0)

	assert.Nil(t, db.AddSubscriptions(dbConn, id, subscriptions))

	subs60, err = db.GetSubscriptionData(dbConn, 60, []string{"testAddr"})
	assert.Nil(t, err)
	assert.Len(t, subs60, 1)
}

func TestDb_DuplicateEntries(t *testing.T) {
	setup.CleanupPgContainer(dbConn)
	var subscriptions []models.SubscriptionData

	id := uint(1)

	for i := 0; i < 10; i++ {
		subscriptions = append(subscriptions, models.SubscriptionData{
			SubscriptionId: id,
			Coin:           60,
			Address:        "testAddr",
		})
	}

	assert.Nil(t, db.AddSubscriptions(dbConn, id, subscriptions))

	subs, err := db.GetSubscriptionData(dbConn, 60, []string{"testAddr"})
	assert.Nil(t, err)
	assert.NotNil(t, subs)
	assert.Equal(t, 1, len(subs))
	assert.True(t, containSub(subscriptions[0], subs))
}

func TestDb_FindSubscriptions_Multiple(t *testing.T) {
	setup.CleanupPgContainer(dbConn)

	var subscriptions []models.SubscriptionData
	subscriptions = append(subscriptions, models.SubscriptionData{
		Coin:    60,
		Address: "testAddr",
	})

	for i := 1; i < 6; i++ {
		subscriptions[0].SubscriptionId = uint(i)
		assert.Nil(t, db.AddSubscriptions(dbConn, uint(i), subscriptions))
	}

	subscriptions[0].SubscriptionId = uint(1)
	assert.Nil(t, db.AddSubscriptions(dbConn, uint(1), subscriptions))

	subs, err := db.GetSubscriptionData(dbConn, 60, []string{"testAddr"})
	assert.Nil(t, err)
	assert.Equal(t, 5, len(subs))

	for i := 0; i < 5; i++ {
		assert.Equal(t, uint(i)+1, subs[i].SubscriptionId)
	}
}

func TestDb_AddToExisting(t *testing.T) {
	setup.CleanupPgContainer(dbConn)

	var subscriptions []models.SubscriptionData
	subscriptions = append(subscriptions, models.SubscriptionData{
		Coin:    60,
		Address: "testAddr",
	})

	subscriptions[0].SubscriptionId = uint(1)
	assert.Nil(t, db.AddSubscriptions(dbConn, uint(1), subscriptions))

	subs, err := db.GetSubscriptionData(dbConn, 60, []string{"testAddr"})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(subs))

	assert.Equal(t, uint(1), subs[0].SubscriptionId)

	assert.Nil(t, db.AddToExistingSubscription(dbConn, uint(1), subscriptions))

	subs2, err2 := db.GetSubscriptionData(dbConn, 60, []string{"testAddr"})
	assert.Nil(t, err2)
	assert.Equal(t, 1, len(subs2))

	assert.Equal(t, uint(1), subs2[0].SubscriptionId)

	assert.Nil(t, db.AddToExistingSubscription(dbConn, uint(2), subscriptions))
	assert.Nil(t, db.AddToExistingSubscription(dbConn, uint(1), subscriptions))

	for i := 1; i < 2; i++ {
		assert.Nil(t, db.AddToExistingSubscription(dbConn, uint(i), subscriptions))
	}
	assert.NotNil(t, db.AddToExistingSubscription(dbConn, uint(0), subscriptions))

}

func containSub(sub models.SubscriptionData, list []models.SubscriptionData) bool {
	for _, s := range list {
		if sub.Address == s.Address && sub.Coin == s.Coin && sub.SubscriptionId == s.SubscriptionId {
			return true
		}
	}
	return false
}