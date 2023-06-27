package mint

import (
	"fmt"
	repoCommon "github.com/LampardNguyen234/astra-integration-test/common"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/assert"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *MintSuite) RunTest() {
	s.Start()
	for i := 0; i < 5; i++ {
		waitBlock := repoCommon.RandInt64() % 20
		s.WaitForBlock(waitBlock)

		msg := fmt.Sprintf("[TEST %v]", i)

		// current mintInfo
		oldMintINfo, err := s.MintInfo()
		assert.NoError(err, msg)

		expNextInflation := s.NextInflationRate(oldMintINfo.Inflation, oldMintINfo.BondedRatio)
		expNextBlockProvision := expNextInflation.
			MulInt(oldMintINfo.StakingSupply).
			QuoInt(sdk.NewInt(int64(oldMintINfo.Params.InflationParameters.BlocksPerYear))).
			TruncateInt()

		s.WaitUntilBlock(oldMintINfo.Height + 1)
		newMintInfo, err := s.MintInfo()
		assert.NoError(err, msg)

		assert.Compare(newMintInfo.CirculatingSupply, oldMintINfo.CirculatingSupply.Add(expNextBlockProvision), assert.OpGTE)
		assert.Compare(newMintInfo.Inflation, expNextInflation, assert.OpEQ)
		assert.Compare(
			newMintInfo.FoundationBalance,
			oldMintINfo.FoundationBalance.Add(newMintInfo.Params.InflationDistribution.Foundation.MulInt(expNextBlockProvision).TruncateInt()),
			assert.OpGTE,
		)
	}

	s.Finished()
}
