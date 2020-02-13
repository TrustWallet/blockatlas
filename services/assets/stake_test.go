package assets

import (
	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/blockatlas/coin"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"testing"
)

var (
	cosmosCoin = coin.Coin{Handle: "cosmos"}
	validators = []blockatlas.Validator{
		{
			ID:     "test1",
			Status: true,
		},
		{
			ID:     "test2",
			Status: true,
		},
	}
	assets = []AssetValidator{
		{
			ID:          "test1",
			Name:        "Spider",
			Description: "yo",
			Website:     "https://tw.com",
			Status:      ValidatorStatus{Disabled: false},
		},
		{
			ID:          "test2",
			Name:        "Man",
			Description: "lo",
			Website:     "https://tw.com",
			Status:      ValidatorStatus{Disabled: true},
		},
	}
	assetsDisabled = []AssetValidator{
		{
			ID:          "test1",
			Name:        "Spider",
			Description: "yo",
			Website:     "https://tw.com",
			Status:      ValidatorStatus{Disabled: true},
		},
		{
			ID:          "test2",
			Name:        "Man",
			Description: "lo",
			Website:     "https://tw.com",
			Status:      ValidatorStatus{Disabled: true},
		},
	}
	expectedStakeValidator = blockatlas.StakeValidator{
		ID: "test1", Status: true,
		Info: blockatlas.StakeValidatorInfo{
			Name:        "Spider",
			Description: "yo",
			Image:       "/cosmos/validators/assets/test1/logo.png",
			Website:     "https://tw.com",
		},
	}
	expectedStakeValidatorDisabled1 = blockatlas.StakeValidator{
		ID: "test1", Status: false,
		Info: blockatlas.StakeValidatorInfo{
			Name:        "Spider",
			Description: "yo",
			Image:       "/cosmos/validators/assets/test1/logo.png",
			Website:     "https://tw.com",
		},
	}
	expectedStakeValidatorDisabled2 = blockatlas.StakeValidator{
		ID: "test2", Status: false,
		Info: blockatlas.StakeValidatorInfo{
			Name:        "Man",
			Description: "lo",
			Image:       "/cosmos/validators/assets/test2/logo.png",
			Website:     "https://tw.com",
		},
	}
)

func TestGetImage(t *testing.T) {
	assetsURL := "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains"
	image := getImage(cosmosCoin, "TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp", assetsURL)
	expected := "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/cosmos/validators/assets/TGzz8gjYiYRqpfmDwnLxfgPuLVNmpCswVp/logo.png"
	assert.Equal(t, expected, image)
}

func TestCalcAnnual(t *testing.T) {
	type args struct {
		annual     float64
		commission float64
	}

	tests := []struct {
		name   string
		args   args
		wanted float64
	}{
		{
			name: "test TestCalcAnnual 1",
			args: args{
				annual:     10,
				commission: 10,
			},
			wanted: 9,
		},
		{
			name: "test TestCalcAnnual 2",
			args: args{
				annual:     100,
				commission: 10,
			},
			wanted: 90,
		},
		{
			name: "test TestCalcAnnual 3",
			args: args{
				annual:     20,
				commission: 10,
			},
			wanted: 18,
		},
		{
			name: "test TestCalcAnnual 1",
			args: args{
				annual:     30,
				commission: 10,
			},
			wanted: 27,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfo := calculateAnnual(tt.args.annual, tt.args.commission)
			assert.Equal(t, tt.wanted, gotInfo)
		})
	}
}

func TestNormalizeValidator(t *testing.T) {
	result := normalizeValidator(validators[0], assets[0], cosmosCoin)
	assert.Equal(t, expectedStakeValidator, result)
}

func Test_normalizeValidators(t *testing.T) {
	type args struct {
		assets     AssetValidators
		validators []blockatlas.Validator
		coin       coin.Coin
		onlyActive bool
	}
	tests := []struct {
		name string
		args args
		want blockatlas.StakeValidators
	}{
		{"normalize validator", args{assets, validators, cosmosCoin, true}, blockatlas.StakeValidators{expectedStakeValidator}},
		{"normalize active validator", args{assets, validators, cosmosCoin, false}, blockatlas.StakeValidators{expectedStakeValidator, expectedStakeValidatorDisabled2}},
		{"normalize validator with disabled assets", args{assetsDisabled, validators, cosmosCoin, true}, blockatlas.StakeValidators{}},
		{"normalize active validator with disabled assets", args{assetsDisabled, validators, cosmosCoin, false}, blockatlas.StakeValidators{expectedStakeValidatorDisabled1, expectedStakeValidatorDisabled2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeValidators(tt.args.assets, tt.args.validators, tt.args.coin, tt.args.onlyActive)
			assert.Equal(t, tt.want, got)
		})
	}
}
