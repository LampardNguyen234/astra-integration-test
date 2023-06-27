package common

import (
	"context"
	"github.com/LampardNguyen234/astra-go-sdk/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

func (c *TestClient) BalanceCheckInt(address string, expectedAmt sdk.Int, op CompareOP) {
	balance, err := c.Balance(address)
	if err != nil {
		c.Log.Panicf("failed to get balance of %v: %v", address, err)
	}

	if !op.CompareSdkInt(balance.Total, expectedAmt) {
		c.Log.Panicf("expect balance of %v to be %v %v, got %v", address, op.String(), expectedAmt, balance.Total)
	}
}

func (c *TestClient) BalanceCheckFloat64(address string, expectedAmt float64, op CompareOP) {
	balance, err := c.Balance(address)
	if err != nil {
		c.Log.Panicf("failed to get balance of %v: %v", address, err)
	}

	tmp := common.BigIntToFloat64(balance.Total.BigInt())
	if !op.CompareFloat64(tmp, expectedAmt) {
		c.Log.Panicf("expect balance of %v to be %v %v, got %v", address, op.String(), expectedAmt, tmp)
	}
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
