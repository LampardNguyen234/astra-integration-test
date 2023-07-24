package common

import (
	"github.com/LampardNguyen234/astra-integration-test/test-suite/assert"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (c *TestClient) ExpectBalance(address string, expectedAmt sdk.Int, op assert.CompareOP) {
	balance, err := c.Balance(address)
	if err != nil {
		c.Log.Panicf("failed to get balance of %v: %v", address, err)
	}

	assert.Compare(balance.Total, expectedAmt, op)
}
