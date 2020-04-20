package main

import (
	"context"
	"github.com/spf13/viper"
	"github.com/trustwallet/blockatlas/db"
	_ "github.com/trustwallet/blockatlas/docs"
	"github.com/trustwallet/blockatlas/internal"
	"github.com/trustwallet/blockatlas/mq"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/blockatlas/pkg/servicerepo"
	"github.com/trustwallet/blockatlas/services/observer/subscriber"
	"time"
)

const (
	defaultConfigPath = "../../config.yml"
)

var (
	confPath    string
	serviceRepo *servicerepo.ServiceRepo
	database    *db.Instance
)

func init() {
	_, confPath = internal.ParseArgs("", defaultConfigPath)
	serviceRepo = servicerepo.New()

	internal.InitConfig(confPath)
	logger.InitLogger()

	pgUri := viper.GetString("postgres.uri")

	mqHost := viper.GetString("observer.rabbitmq.uri")
	prefetchCount := viper.GetInt("observer.rabbitmq.consumer.prefetch_count")

	internal.InitRabbitMQ(mqHost, prefetchCount)
	subscriber.InitService(serviceRepo)

	var err error
	database, err = db.New(pgUri)
	if err != nil {
		logger.Fatal(err)
	}

	go mq.FatalWorker(time.Second * 10)
	go db.RestoreConnectionWorker(database, time.Second*10, pgUri)
	time.Sleep(time.Millisecond)
}

func main() {
	defer mq.Close()
	if err := mq.Subscriptions.Declare(); err != nil {
		logger.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	subscriberService := subscriber.GetService(serviceRepo)
	go mq.Subscriptions.RunConsumerWithCancelAndDbConn(subscriberService.RunSubscriber, database, ctx)

	internal.SetupGracefulShutdownForObserver(cancel)
}
