// +build ignore

//go:generate rm -f coins.go
//go:generate go run gen.go

package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"text/template"
	"time"
)

const (
	coinFile = "../coins.yml"
	filename = "coins.go"
)

type Coin struct {
	ID               uint   `yaml:"id"`
	Handle           string `yaml:"handle"`
	Symbol           string `yaml:"symbol"`
	Title            string `yaml:"name"`
	Decimals         uint   `yaml:"decimals"`
	BlockTime        int    `yaml:"blockTime"`
	MinConfirmations int64  `yaml:"minConfirmations"`
	SampleAddr       string `yaml:"sampleAddress"`
}

func main() {
	var coinList []Coin
	coin, err := os.Open(coinFile)
	dec := yaml.NewDecoder(coin)
	err = dec.Decode(&coinList)
	if err != nil {
		log.Panic(err)
	}

	f, err := os.Create(filename)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	err = coinsTemplate.Execute(f, map[string]interface{}{
		"Timestamp": time.Now(),
		"Coins":     coinList,
	})
	if err != nil {
		log.Panic(err)
	}
}

var coinsTemplate = template.Must(template.New("").Parse(
	`// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}
// using data from coins.yml
package coin

import (
	"fmt"
)

// Coin is the native currency of a blockchain
type Coin struct {
	ID               uint
	Handle           string
	Symbol           string
	Title            string
	Decimals         uint
	BlockTime        int
	MinConfirmations int64
	SampleAddr       string
}

func (c Coin) String() string {
	return fmt.Sprintf("[%s] %s (#%d)", c.Symbol, c.Title, c.ID)
}

const (
{{- range .Coins }}
	{{ .Symbol }} = {{ .ID }}
{{- end }}
)

var Coins = map[uint]Coin{
{{- range .Coins }}
	{{.ID}}: {
		ID:               {{.ID}},
		Handle:           "{{.Handle}}",
		Symbol:           "{{.Symbol}}",
		Title:            "{{.Title}}",
		Decimals:         {{.Decimals}},
		BlockTime:        {{.BlockTime}},
		MinConfirmations: {{.MinConfirmations}},
		SampleAddr:       "{{.SampleAddr}}",
	},
{{- end }}
}
`))
