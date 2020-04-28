package tezos

import (
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"github.com/trustwallet/blockatlas/services/assets"
)

const (
	Annual             = 6.09
	LockTime           = 0
	MinimumStakeAmount = "0"
)

func (p *Platform) GetDelegations(address string) (blockatlas.DelegationsPage, error) {
	account, err := p.rpcClient.GetAccount(address)
	if err != nil {
		return nil, err
	}
	if len(account.Delegate) == 0 {
		return make(blockatlas.DelegationsPage, 0), nil
	}

	validators, err := assets.GetValidatorsMap(p)
	if err != nil {
		return nil, err
	}
	return NormalizeDelegation(account, validators)
}

func NormalizeDelegation(account Account, validators blockatlas.ValidatorMap) (blockatlas.DelegationsPage, error) {
	validator, ok := validators[account.Delegate]
	if !ok {
		logger.Warn("Validator not found", logger.Params{"address": validator.ID, "platform": "tezos", "delegation": account.Delegate})
		validator = getUnknownValidator(validator.ID)
	}
	return blockatlas.DelegationsPage{
		{
			Delegator: validator,
			Value:     account.Balance,
			Status:    blockatlas.DelegationStatusActive,
		},
	}, nil
}

func (p *Platform) GetValidators() (blockatlas.ValidatorPage, error) {
	results := make(blockatlas.ValidatorPage, 0)

	validators, err := p.getCurrentValidators()
	if err != nil {
		return results, err
	}

	for _, v := range validators {
		results = append(results, normalizeValidator(v))
	}
	return results, nil
}

func (p *Platform) getCurrentValidators() (validators []Validator, err error) {
	periodType, err := p.rpcClient.GetPeriodType()
	if err != nil {
		return validators, err
	}

	switch periodType {
	case TestingPeriodType:
		return p.rpcClient.GetValidators("head~32768")
	default:
		return p.rpcClient.GetValidators("head")
	}
}

func (p *Platform) GetDetails() blockatlas.StakingDetails {
	return getDetails()
}

func (p *Platform) UndelegatedBalance(address string) (string, error) {
	account, err := p.rpcClient.GetAccount(address)
	if err != nil {
		return "0", err
	}
	return account.Balance, nil
}

func getDetails() blockatlas.StakingDetails {
	return blockatlas.StakingDetails{
		Reward:        blockatlas.StakingReward{Annual: Annual},
		MinimumAmount: "0",
		LockTime:      0,
		Type:          blockatlas.DelegationTypeDelegate,
	}
}

func normalizeValidator(v Validator) (validator blockatlas.Validator) {
	// How to calculate Tezos APR? I have no idea. Tezos team does not know either. let's assume it's around 7% - no way to calculate in decentralized manner
	// Delegation rewards distributed by the validators manually, it's up to them to do it.
	return blockatlas.Validator{
		Status:  true,
		ID:      v.Address,
		Details: getDetails(),
	}
}

func getUnknownValidator(address string) blockatlas.StakeValidator {
	return blockatlas.StakeValidator{
		ID:     address,
		Status: false,
		Info: blockatlas.StakeValidatorInfo{
			Name:        "Decommissioned",
			Description: "Decommissioned",
		},
		Details: blockatlas.StakingDetails{
			Reward: blockatlas.StakingReward{
				Annual: Annual,
			},
			LockTime:      LockTime,
			MinimumAmount: MinimumStakeAmount,
			Type:          blockatlas.DelegationTypeDelegate,
		},
	}
}
