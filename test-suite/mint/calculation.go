package mint

import (
	"github.com/AstraProtocol/astra/v2/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *MintSuite) NextInflationRate(inflation, bondedRatio sdk.Dec) sdk.Dec {
	// (1 - bondedRatio/GoalBonded) * InflationRateChange
	inflationRateChangePerYear := sdk.OneDec().
		Sub(bondedRatio.Quo(s.params.InflationParameters.GoalBonded)).
		Mul(s.params.InflationParameters.InflationRateChange)

	inflationRateChange := inflationRateChangePerYear.QuoInt(sdk.NewIntFromUint64(s.params.InflationParameters.BlocksPerYear))
	// adjust the new annual inflation for this next cycle
	inflation = inflation.Add(inflationRateChange) // note inflationRateChange may be negative
	if inflation.GT(s.params.InflationParameters.InflationMax) {
		inflation = s.params.InflationParameters.InflationMax
	}
	if inflation.LT(s.params.InflationParameters.InflationMin) {
		inflation = s.params.InflationParameters.InflationMin
	}

	return inflation
}

// getProportion returns the amount of coins given its distribution and the total amount of coins.
func getProportion(
	amt sdk.Int,
	distribution sdk.Dec,
) sdk.Int {
	return sdk.NewDecFromInt(amt).Mul(distribution).TruncateInt()
}

func getAllProportions(blkProvision sdk.Int, distribution types.InflationDistribution) (
	staking sdk.Int, foundation sdk.Int, community sdk.Int) {
	staking = getProportion(blkProvision, distribution.StakingRewards)
	foundation = getProportion(blkProvision, distribution.Foundation)
	community = blkProvision.Sub(foundation).Sub(staking)
	return
}
