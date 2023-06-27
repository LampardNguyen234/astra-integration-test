package mint

import (
	"fmt"
	repoCommon "github.com/LampardNguyen234/astra-integration-test/common"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *MintSuite) RunTest() {
	numTests := 10
	for i := 0; i < numTests; i++ {
		waitBlock := repoCommon.RandInt64() % 20
		s.WaitForBlock(waitBlock)

		msg := fmt.Sprintf("[TEST %v]", i)

		// current mintInfo
		oldMintINfo, err := s.MintInfo()
		common.NoError(err, msg)

		expNextInflation := s.NextInflationRate(oldMintINfo.Inflation, oldMintINfo.BondedRatio)
		expNextBlockProvision := expNextInflation.
			MulInt(oldMintINfo.StakingSupply).
			QuoInt(sdk.NewInt(int64(oldMintINfo.Params.InflationParameters.BlocksPerYear))).
			TruncateInt()

		s.WaitUntilBlock(oldMintINfo.Height + 1)
		newMintInfo, err := s.MintInfo()
		common.NoError(err, msg)

		s.Compare(newMintInfo.CirculatingSupply, oldMintINfo.CirculatingSupply.Add(expNextBlockProvision), common.OpGTE)
		s.Compare(newMintInfo.Inflation, expNextInflation, common.OpEQ)
		s.Compare(
			newMintInfo.FoundationBalance,
			oldMintINfo.FoundationBalance.Add(newMintInfo.Params.InflationDistribution.Foundation.MulInt(expNextBlockProvision).TruncateInt()),
			common.OpGTE,
		)

		s.Log.Debugf("TEST %v PASS", i)
	}
}
