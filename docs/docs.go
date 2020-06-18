// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/ns/lookup": {
            "get": {
                "description": "Lookup ENS/ZNS to find registered addresses",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Naming"
                ],
                "summary": "Lookup .eth / .zil addresses",
                "operationId": "lookup",
                "parameters": [
                    {
                        "type": "string",
                        "description": "string name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "string coin",
                        "name": "coin",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/blockatlas.Resolved"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/endpoint.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v2/ns/lookup": {
            "get": {
                "description": "Lookup ENS/ZNS to find registered addresses for multiple coins",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Naming"
                ],
                "summary": "Lookup .eth / .zil addresses",
                "operationId": "lookup",
                "parameters": [
                    {
                        "type": "string",
                        "description": "string name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "List of coins",
                        "name": "coins",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/blockatlas.Resolved"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/endpoint.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v2/staking/delegations": {
            "post": {
                "description": "Get Stake Delegations for multiple coins",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Staking"
                ],
                "summary": "Get Multiple Stake Delegations",
                "operationId": "batch_delegations",
                "parameters": [
                    {
                        "description": "Validators addresses and coins",
                        "name": "delegations",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/endpoint.AddressesRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/blockatlas.DelegationsBatchPage"
                        }
                    }
                }
            }
        },
        "/v2/staking/list": {
            "post": {
                "description": "Get Stake Delegations for multiple coins",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Staking"
                ],
                "summary": "Get Multiple Stake Delegations",
                "operationId": "batch_delegations",
                "parameters": [
                    {
                        "description": "Validators addresses and coins",
                        "name": "delegations",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/endpoint.AddressesRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/blockatlas.DelegationsBatchPage"
                        }
                    }
                }
            }
        },
        "/v2/tokens": {
            "post": {
                "description": "Get tokens",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transactions"
                ],
                "summary": "Get list of tokens by map: coin -\u003e [addresses]",
                "operationId": "tokens_v3",
                "parameters": [
                    {
                        "default": "{\"60\": [\"0xb3624367b1ab37daef42e1a3a2ced012359659b0\"]}",
                        "description": "Payload",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/blockatlas.ResultsResponse"
                        }
                    }
                }
            }
        },
        "/v2/{coin}/staking/delegations/{address}": {
            "get": {
                "description": "Get stake delegations from the address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Staking"
                ],
                "summary": "Get Stake Delegations",
                "operationId": "delegations",
                "parameters": [
                    {
                        "type": "string",
                        "default": "tron",
                        "description": "the coin name",
                        "name": "coin",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "TPJYCz8ppZNyvw7pTwmjajcx4Kk1MmEUhD",
                        "description": "the query address",
                        "name": "address",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/blockatlas.DelegationResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/endpoint.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v2/{coin}/staking/validators": {
            "get": {
                "description": "Get validators from the address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Staking"
                ],
                "summary": "Get Validators",
                "operationId": "validators",
                "parameters": [
                    {
                        "type": "string",
                        "default": "cosmos",
                        "description": "the coin name",
                        "name": "coin",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/blockatlas.DocsResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/endpoint.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v2/{coin}/tokens/{address}": {
            "get": {
                "description": "Get tokens from the address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transactions"
                ],
                "summary": "Get Tokens",
                "operationId": "tokens",
                "parameters": [
                    {
                        "type": "string",
                        "default": "ethereum",
                        "description": "the coin name",
                        "name": "coin",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "0x5574Cd97432cEd0D7Caf58ac3c4fEDB2061C98fB",
                        "description": "the query address",
                        "name": "address",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/blockatlas.CollectionPage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/endpoint.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v2/{coin}/transactions/xpub/{xpub}": {
            "get": {
                "description": "Get transactions from XPUB address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transactions"
                ],
                "summary": "Get Transactions by XPUB",
                "operationId": "tx_xpub_v2",
                "parameters": [
                    {
                        "type": "string",
                        "default": "bitcoin",
                        "description": "the coin name",
                        "name": "coin",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "zpub6ruK9k6YGm8BRHWvTiQcrEPnFkuRDJhR7mPYzV2LDvjpLa5CuGgrhCYVZjMGcLcFqv9b2WvsFtY2Gb3xq8NVq8qhk9veozrA2W9QaWtihrC",
                        "description": "the xpub key",
                        "name": "xpub",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/endpoint.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v2/{coin}/transactions/{address}": {
            "get": {
                "description": "Get transactions from the address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transactions"
                ],
                "summary": "Get Transactions",
                "operationId": "tx_v2",
                "parameters": [
                    {
                        "type": "string",
                        "default": "tezos",
                        "description": "the coin name",
                        "name": "coin",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "tz1WCd2jm4uSt4vntk4vSuUWoZQGhLcDuR9q",
                        "description": "the query address",
                        "name": "address",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/endpoint.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v3/staking/list": {
            "get": {
                "description": "Get staking info by coin ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Staking"
                ],
                "summary": "Get staking info by coin ID",
                "operationId": "batch_info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "List of coins",
                        "name": "coins",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/blockatlas.DelegationsBatchPage"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/endpoint.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/v4/collectibles/categories": {
            "post": {
                "description": "Get collection categories",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Collections"
                ],
                "summary": "Get list of collections from a specific coin and addresses",
                "operationId": "collection_categories_v4",
                "parameters": [
                    {
                        "default": "{\"60\": [\"0xb3624367b1ab37daef42e1a3a2ced012359659b0\"]}",
                        "description": "Payload",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/blockatlas.DocsResponse"
                        }
                    }
                }
            }
        },
        "/v4/{coin}/collections/{owner}/collection/{collection_id}": {
            "get": {
                "description": "Get a collection from the address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Collections"
                ],
                "summary": "Get Collection",
                "operationId": "collection_v4",
                "parameters": [
                    {
                        "type": "string",
                        "default": "ethereum",
                        "description": "the coin name",
                        "name": "coin",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "0x0875BCab22dE3d02402bc38aEe4104e1239374a7",
                        "description": "the query address",
                        "name": "owner",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "0x06012c8cf97bead5deae237070f9587f8e7a266d",
                        "description": "the query collection",
                        "name": "collection_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/blockatlas.CollectionPage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/endpoint.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "blockatlas.Collection": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "coin": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "external_link": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "image_url": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "blockatlas.CollectionPage": {
            "type": "array",
            "items": {
                "$ref": "#/definitions/blockatlas.Collection"
            }
        },
        "blockatlas.Delegation": {
            "type": "object",
            "properties": {
                "delegator": {
                    "type": "object",
                    "$ref": "#/definitions/blockatlas.StakeValidator"
                },
                "metadata": {
                    "type": "object"
                },
                "status": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "blockatlas.DelegationResponse": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "balance": {
                    "type": "string"
                },
                "coin": {
                    "type": "object",
                    "$ref": "#/definitions/coin.ExternalCoin"
                },
                "delegations": {
                    "type": "object",
                    "$ref": "#/definitions/blockatlas.DelegationsPage"
                },
                "details": {
                    "type": "object",
                    "$ref": "#/definitions/blockatlas.StakingDetails"
                }
            }
        },
        "blockatlas.DelegationsBatchPage": {
            "type": "array",
            "items": {
                "$ref": "#/definitions/blockatlas.DelegationResponse"
            }
        },
        "blockatlas.DelegationsPage": {
            "type": "array",
            "items": {
                "$ref": "#/definitions/blockatlas.Delegation"
            }
        },
        "blockatlas.DocsResponse": {
            "type": "object",
            "properties": {
                "docs": {
                    "type": "object"
                }
            }
        },
        "blockatlas.Resolved": {
            "type": "object",
            "properties": {
                "coin": {
                    "type": "integer"
                },
                "result": {
                    "type": "string"
                }
            }
        },
        "blockatlas.ResultsResponse": {
            "type": "object",
            "properties": {
                "docs": {
                    "type": "object"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "blockatlas.StakeValidator": {
            "type": "object",
            "properties": {
                "details": {
                    "type": "object",
                    "$ref": "#/definitions/blockatlas.StakingDetails"
                },
                "id": {
                    "type": "string"
                },
                "info": {
                    "type": "object",
                    "$ref": "#/definitions/blockatlas.StakeValidatorInfo"
                },
                "status": {
                    "type": "boolean"
                }
            }
        },
        "blockatlas.StakeValidatorInfo": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "image": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "website": {
                    "type": "string"
                }
            }
        },
        "blockatlas.StakingDetails": {
            "type": "object",
            "properties": {
                "locktime": {
                    "type": "integer"
                },
                "minimum_amount": {
                    "type": "string"
                },
                "reward": {
                    "type": "object",
                    "$ref": "#/definitions/blockatlas.StakingReward"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "blockatlas.StakingReward": {
            "type": "object",
            "properties": {
                "annual": {
                    "type": "number"
                }
            }
        },
        "coin.ExternalCoin": {
            "type": "object",
            "properties": {
                "coin": {
                    "type": "integer"
                },
                "decimals": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "symbol": {
                    "type": "string"
                }
            }
        },
        "endpoint.AddressBatchRequest": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "coin": {
                    "type": "integer"
                }
            }
        },
        "endpoint.AddressesRequest": {
            "type": "array",
            "items": {
                "$ref": "#/definitions/endpoint.AddressBatchRequest"
            }
        },
        "endpoint.ErrorDetails": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "endpoint.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "object",
                    "$ref": "#/definitions/endpoint.ErrorDetails"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
