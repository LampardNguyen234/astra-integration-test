package common

import (
	"github.com/LampardNguyen234/astra-go-sdk/common"
)

func (c *TestClient) BalanceCheck(address string, expectedAmt float64, op CompareOP) {
	balance, err := c.Balance(address)
	if err != nil {
		c.Log.Panicf("failed to get balance of %v: %v", address, err)
	}

	tmp := common.BigIntToFloat64(balance.Total.BigInt())
	if !op.Compare(tmp, expectedAmt) {
		c.Log.Panicf("expect balance of %v to be %v %v, got %v", address, op.String(), expectedAmt, tmp)
	}
}
