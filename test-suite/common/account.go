package common

import (
	"context"
	"github.com/LampardNguyen234/astra-go-sdk/common"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/assert"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

func (c *TestClient) BalanceCompare(address string, expectedAmt sdk.Int, op assert.CompareOP) {
	balance, err := c.Balance(address)
	if err != nil {
		c.Log.Panicf("failed to get balance of %v: %v", address, err)
	}

	assert.Compare(balance.Total, expectedAmt, op)
}

func (c *TestClient) WaitForBalanceUpdated(addr string, expectedAmt float64) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	amt := common.Float64ToBigInt(expectedAmt)
	for {
		select {
		case <-ctx.Done():
			c.Log.Panicf("failed to check balance of addr %v: TIMED-OUT", addr)
		default:
			time.Sleep(2 * time.Second)
			balance, _ := c.Balance(addr)
			if balance != nil {
				if balance.Total.GTE(sdk.NewIntFromBigInt(amt)) {
					return
				}
			}
		}
	}
}
