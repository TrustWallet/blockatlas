# Block Atlas by Trust Wallet

![Go Version](https://img.shields.io/github/go-mod/go-version/TrustWallet/blockatlas)
![CI](https://github.com/trustwallet/blockatlas/workflows/CI/badge.svg)
[![codecov](https://codecov.io/gh/trustwallet/blockatlas/branch/master/graph/badge.svg)](https://codecov.io/gh/trustwallet/blockatlas)
[![Go Report Card](https://goreportcard.com/badge/trustwallet/blockatlas)](https://goreportcard.com/report/TrustWallet/blockatlas)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=trustwallet/blockatlas)](https://dependabot.com)

> BlockAtlas is a clean explorer API and transaction observer for cryptocurrencies.

BlockAtlas connects to nodes or explorer APIs of the supported coins and maps transaction data,
account transaction history into a generic, easy to work with JSON format.
It is in production use at the [Trust Wallet app](https://trustwallet.com/), 
the official cryptocurrency wallet of Binance. Also is in production at the [BUTTON Wallet](https://buttonwallet.com), Telegram based non-custodial wallet.
The observer API watches the chain for new transactions and generates notifications by guids.

#### Supported Coins

Block Atlas supports more than 25 blockchains: Bitcoin, Ethereum, Binance Chain etc, The full feature matrix is [here](docs/features.csv).

## Architecture

#### NOTE

Currently Block Atlas is under active development and is not well documented. If you still want to run it on your own or help to contribute, **please** pay attention that currently integration, newman, functional tests are not working locally without all endpoints. We are fixing that issue and soon you will be able to test all the stuff locally

Blockatlas allows to:

-   Get information about transactions, tokens, staking details, collectibles for supported coins.
-   Subscribe for price notifications via Rabbit MQ

Platform API is independent service and can work with the specific blockchain only (like Bitcoin, Ethereum, etc)

Notifications:

-   Subscriber Producer - Create new blockatlas.SubscriptionEvent [Not implemented at Atlas, write it on your own]

-   Subscriber - Get subscriptions from queue, set them to the DB

-   Parser - Parse the block, convert block to the transactions batch, send to queue

-   Notifier - Check each transaction for having the same address as stored at DB, if so - send tx data and id to the next queue

-   Notifier Consumer - Notify the user [Not implemented at Atlas, write it on your own]


```
New Subscriptions --(Rabbit MQ)--> Subscriber --> DB
                                                   |
                      Parser  --(Rabbit MQ)--> Notifier --(Rabbit MQ)--> Notifier Consumer --> User

```

The whole flow is not available at Atlas repo. We will have integration tests with it. Also there will be examples of all instances soon.

## Setup

### Prerequisite

-   [Go Toolchain](https://golang.org/doc/install) versions 1.14+

    Depends on what type of Blockatlas service you would like to run will also be needed.
-   [Postgres](https://www.postgresql.org/download) to store user subscriptions and latest parsed block number
-   [Rabbit MQ](https://www.rabbitmq.com/#getstarted) to pass subscriptions and send transaction notifications

### Quick Start

#### Get source code

Download source to `GOPATH`

```shell
go get -u github.com/trustwallet/blockatlas
cd $(go env GOPATH)/src/github.com/trustwallet/blockatlas
```

#### Build and run

Read [configuration](#configuration) info

```shell
# Start Platform API server at port 8420 with the path to the config.yml ./
go build -o api-bin cmd/api/main.go && ./api-bin -p 8420

# Start parser with the path to the config.yml ./ 
go build -o parser-bin cmd/parser/main.go && ./parser-bin

# Start notifier with the path to the config.yml ./ 
go build -o notifier-bin cmd/notifier/main.go && ./notifier-bin

# Start subscriber with the path to the config.yml ./ 
go build -o subscriber-bin cmd/subscriber/main.go && ./subscriber-bin
```

### make command

Build and start all services:

```shell
make go-build
make start
```

Build and start individual service:

```shell
make go-build-api
make start
```

### Docker

Build and run all services:

```shell
docker-compose build
docker-compose up
```

Build and run individual service:

```shell
docker-compose build api
docker-compose start api
```

## Configuration

When any of Block Atlas services started they look up inside [default configuration](./config.yml).
Most coins offering public RPC/explorer APIs are enabled, thus Block Atlas can be started and used right away, no additional configuration needed.
By default starting any of the [services](#architecture) will enable all platforms

To run a specific service only by passing environmental variable, e.g: `ATLAS_PLATFORM=ethereum` :

```shell
ATLAS_PLATFORM=ethereum go run cmd/api/main.go

ATLAS_PLATFORM=ethereum binance bitcoin go run cmd/api/main.go # for multiple platforms
```

or change in config file

```yaml
# Single
platform: [ethereum]
# Multiple 
platform: [ethereum, binance, bitcoin]
```

This way you can one platform per binary, for scalability and sustainability.

To enable use of private endpoint:

```yaml
nimiq:
  api: http://localhost:8648
```

It works the same for worker - you can run all observer at 1 binary or 30 coins per 30 binaries

#### Environment

The rest gets loaded from environment variables.
Every config option is available under the `ATLAS_` prefix. Nested keys are joined via `_`.

Example:

```shell
ATLAS_NIMIQ_API=http://localhost:8648
```

## Tests

### Unit tests

    make test

### Mocked tests

End-to-end tests with calls to external APIs has great value, but they are not suitable for regular CI verification, beacuse any external reason could break the tests.

    # Start API server with mocked config, at port 8437 ./ 
    go build -o api-bin cmd/api/main.go && ./api-bin -p 8437 -c configmock.yml

Therefore mocked API-level tests are used, whereby external APIs are replaced by mocks.

-   External mocks are implemented as a simple, own, golang `mockserver`.  It listens locally, and returns responses to specific API paths, taken from json data files.
-   There is a file where API paths and corresponding data files are listed.
-   Tests invoke into blockatlas through public APIs only, and are executed using _newman_ (Postman cli -- `make newman-mocked`).
-   Product code, and even test code should not be aware whether it runs with mocks or the real external endpoints.
-   See Makefile for targets with 'mock'; platform can be started locally with mocks using `make start-platform-api-mock`.
-   The newman tests can be executed with unmocked external APIs as well, but verifications may fail, because some APIs return variable responses.  Unmocked tests are not intended for regular CI execution, but as ad-hoc development tests.
-   General steps for creating new mocked tests: replace endpoint to localhost:3347, observe incoming calls (visible in mockserver's output), obtain real response from external API (with exact same parameters), place response in a file, add path + file to data file list.  Restart mock, and verify that blockatlas provides correct output.  Also, add verifications of results to the tests.

## Docs

Swagger API docs provided at path `/swagger/index.html`

or you can install `go-swagger` and render it locally (macOS example)

Install:

```shell
brew tap go-swagger/go-swagger
brew install go-swagger
```

Render: 

```shell
swagger serve docs/swagger.yaml
```

#### Updating Docs

-   After creating a new route, add comments to your API source code, [See Declarative Comments Format](https://swaggo.github.io/swaggo.io/declarative_comments_format/).

-   Run `$ make go-gen-docs` in root folder.

## Contributing

If you'd like to add support for a new blockchain, feel free to file a pull request.
Note that most tokens that run on top of other chains are already supported and
don't require code changes (e.g. ERC-20).

The best way to submit feedback and report bugs is to open a GitHub issue.
Please be sure to include your operating system, version number, and
[steps](https://gist.github.com/nrollr/eb24336b8fb8e7ba5630) to reproduce reported bugs.

More resources for developers are in [CONTRIBUTING.md](CONTRIBUTING.md).
