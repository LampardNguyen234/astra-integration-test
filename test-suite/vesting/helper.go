package vesting

import (
	"github.com/LampardNguyen234/astra-go-sdk/client/msg_params"
	sdkCommon "github.com/LampardNguyen234/astra-go-sdk/common"
	"github.com/LampardNguyen234/astra-integration-test/common"
	vestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

func randPeriods(totalAmt float64) vestingTypes.Periods {
	vestingLength := 1 + common.RandUint64()%20
	duration := int64(vestingLength) + int64(common.RandUint64()%30)
	tmp := &msg_params.TxCreateVestingParams{
		VestingDuration: duration,
		Amount:          sdkCommon.Float64ToBigInt(totalAmt),
		VestingLength:   uint(vestingLength),
	}
	ret := tmp.VestingPeriods()
	for _, p := range ret {
		if p.Length <= 0 {
			p.Length = 1
		}
	}

	return ret
}

func randPeriodsWithDuration(totalAmt float64, duration int64) vestingTypes.Periods {
	tmp := &msg_params.TxCreateVestingParams{
		VestingDuration: duration,
		Amount:          sdkCommon.Float64ToBigInt(totalAmt),
		VestingLength:   uint(1 + common.RandUint64()%10),
	}

	return tmp.VestingPeriods()
}
