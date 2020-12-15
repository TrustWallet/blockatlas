package main

import (
	"context"
	"time"

	"github.com/trustwallet/golibs/network/middleware"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/blockatlas/api"
	"github.com/trustwallet/blockatlas/config"
	"github.com/trustwallet/blockatlas/db"
	_ "github.com/trustwallet/blockatlas/docs"
	"github.com/trustwallet/blockatlas/internal"
	"github.com/trustwallet/blockatlas/mq"
	"github.com/trustwallet/blockatlas/platform"
	"github.com/trustwallet/blockatlas/services/tokenindexer"
	"github.com/trustwallet/blockatlas/services/tokensearcher"
)

const (
	defaultPort       = "8420"
	defaultConfigPath = "../../config.yml"
)

var (
	ctx            context.Context
	cancel         context.CancelFunc
	port, confPath string
	engine         *gin.Engine
	database       *db.Instance
	ts             tokensearcher.Instance
	ti             tokenindexer.Instance
)

func init() {
	port, confPath = internal.ParseArgs(defaultPort, defaultConfigPath)
	ctx, cancel = context.WithCancel(context.Background())
	var err error

	internal.InitConfig(confPath)

	err = middleware.SetupSentry(config.Default.Sentry.DSN)
	if err != nil {
		log.Error(err)
	}

	engine = internal.InitEngine(config.Default.Gin.Mode)
	platform.Init(config.Default.Platform)

	database, err = db.New(config.Default.Postgres.URL, config.Default.Postgres.Log)
	if err != nil {
		log.Fatal(err)
	}
	go database.RestoreConnectionWorker(time.Second*10, config.Default.Postgres.URL)

	internal.InitRabbitMQ(
		config.Default.Observer.Rabbitmq.URL,
		config.Default.Observer.Rabbitmq.Consumer.PrefetchCount,
	)

	if err := mq.TokensRegistration.Declare(); err != nil {
		log.Fatal(err)
	}
	if err := mq.RawTransactionsTokenIndexer.Declare(); err != nil {
		log.Fatal(err)
	}

	ts = tokensearcher.Init(database, platform.TokensAPIs, mq.TokensRegistration)
	ti = tokenindexer.Init(database)

	go mq.FatalWorker(time.Second * 10)
}

func main() {
	api.SetupTokensIndexAPI(engine, ti)
	api.SetupTokensSearcherAPI(engine, ts)
	api.SetupSwaggerAPI(engine)
	api.SetupPlatformAPI(engine)

	internal.SetupGracefulShutdown(ctx, port, engine)
	cancel()
}
