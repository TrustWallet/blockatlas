package main

import (
	"context"
	"time"

	"github.com/trustwallet/blockatlas/services/notifier"

	"github.com/trustwallet/blockatlas/config"
	"github.com/trustwallet/blockatlas/services/subscriber"
	"github.com/trustwallet/blockatlas/services/tokenindexer"
	"github.com/trustwallet/blockatlas/services/tokensearcher"

	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/blockatlas/db"
	"github.com/trustwallet/blockatlas/internal"
	"github.com/trustwallet/blockatlas/mq"
)

const (
	defaultConfigPath = "../../config.yml"
)

var (
	ctx      context.Context
	cancel   context.CancelFunc
	database *db.Instance
)

func init() {
	ctx, cancel = context.WithCancel(context.Background())
	_, confPath := internal.ParseArgs("", defaultConfigPath)

	internal.InitConfig(confPath)

	internal.InitRabbitMQ(
		config.Default.Observer.Rabbitmq.URL,
		config.Default.Observer.Rabbitmq.Consumer.PrefetchCount,
	)

	var err error
	database, err = db.New(config.Default.Postgres.URL, config.Default.Postgres.Read.URL,
		config.Default.Postgres.Log)
	if err != nil {
		log.Fatal("Postgres init: ", err)
	}
	go database.RestoreConnectionWorker(ctx, time.Second*10, config.Default.Postgres.URL)

	time.Sleep(time.Millisecond)
}

func main() {
	defer mq.Close()

	queues := []mq.Queue{
		mq.TxNotifications,
		mq.RawTransactionsTokenIndexer,
		mq.RawTransactions,
		mq.RawTransactionsSearcher,
		mq.Subscriptions,
		mq.TokensRegistration,
	}

	for _, queue := range queues {
		err := queue.Declare()
		if err != nil {
			log.Fatal("Declare ", queue, err)
		}
	}

	go mq.RawTransactions.RunConsumerWithCancelAndDbConn(notifier.RunNotifier, database, ctx)

	go mq.RawTransactionsTokenIndexer.RunConsumerWithCancelAndDbConn(tokenindexer.RunTokenIndexer, database, ctx)
	go mq.RawTransactionsSearcher.RunConsumerWithCancelAndDbConn(tokensearcher.Run, database, ctx)

	go mq.Subscriptions.RunConsumerWithCancelAndDbConn(subscriber.RunTransactionsSubscriber, database, ctx)
	go mq.TokensRegistration.RunConsumerWithCancelAndDbConn(subscriber.RunTokensSubscriber, database, ctx)

	go mq.FatalWorker(time.Second * 10)

	internal.SetupGracefulShutdownForObserver()
	cancel()
}
