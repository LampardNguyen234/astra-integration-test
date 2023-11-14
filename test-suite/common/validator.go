package common

import (
	"github.com/LampardNguyen234/astra-integration-test/common"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (c *TestClient) MustRandActiveValidator() stakingTypes.Validator {
	allValidators, err := c.AllValidators(stakingTypes.Bonded)
	if err != nil {
		c.Log.Panicf("failed to get validator: %v", err)
	}

	return allValidators[int(common.RandUint64()%uint64(len(allValidators)))]
}
