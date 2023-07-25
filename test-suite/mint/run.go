package mint

import (
	"github.com/LampardNguyen234/astra-go-sdk/client"
	. "github.com/LampardNguyen234/astra-integration-test/framework"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func (s *MintSuite) RunTest() {
	s.Start()

	tc := s.registerTests()
	tc.Run()
	tc.Report()

	s.Finished()
}

func (s *MintSuite) registerTests() ITestNode {
	RegisterFailHandler(ginkgo.Fail)

	var totalFee, burnedFee sdk.Int
	var totalStaked, totalUnStaked sdk.Int
	var oldMintInfo, newMintInfo *client.ProvisionInfo
	var err error
	var expNextInflation sdk.Dec
	var expNextBlockProvision, expStakingRewards, expFoundation, expCommunity sdk.Int

	root := Describe(s.Name(),
		When("there are no staking-related transactions",
			Before(func() {
				oldMintInfo, newMintInfo, totalFee, err = s.mintInfoWithNoStakingTxs()
				Expect(err).To(BeNil())

				burnedFee = sdk.MustNewDecFromStr("0.5").MulInt(totalFee).TruncateInt()

				expNextInflation = s.NextInflationRate(oldMintInfo.Inflation, oldMintInfo.BondedRatio)
				expNextBlockProvision = expNextInflation.
					MulInt(oldMintInfo.CirculatingSupply).
					QuoInt(sdk.NewInt(int64(oldMintInfo.Params.InflationParameters.BlocksPerYear))).
					TruncateInt()

				expStakingRewards, expFoundation, expCommunity = getAllProportions(expNextBlockProvision, oldMintInfo.Params.InflationDistribution)
				_ = expStakingRewards
			}),

			It("bonded must stay the same", func() {
				Expect(newMintInfo.StakingSupply.Sub(oldMintInfo.StakingSupply)).To(
					common.EQ(sdk.ZeroInt()),
				)
			}),

			It("circulatingSupply must be correct", func() {
				Expect(newMintInfo.CirculatingSupply.Sub(oldMintInfo.CirculatingSupply)).To(
					Equal(expNextBlockProvision.Sub(burnedFee)))
			}),

			It("inflation must be correct", func() {
				Expect(newMintInfo.Inflation).To(common.LTE(s.params.InflationParameters.InflationMax))
				Expect(newMintInfo.Inflation).To(common.GTE(s.params.InflationParameters.InflationMin))
				Expect(newMintInfo.Inflation).To(Equal(expNextInflation))
			}),

			It("totalMintedProvision must properly increase", func() {
				Expect(newMintInfo.TotalMintedProvision.Sub(oldMintInfo.TotalMintedProvision)).To(Equal(
					expNextBlockProvision))
			}),

			It("foundation balance must be correct", func() {
				Expect(newMintInfo.FoundationBalance.Sub(oldMintInfo.FoundationBalance)).To(Equal(
					expFoundation,
				))
			}),

			It("community balance must be correct", func() {
				// we allow a variant of 1000.
				Expect(newMintInfo.CommunityBalance.Sub(oldMintInfo.CommunityBalance).Sub(expCommunity).Abs()).
					To(common.LTE(
						sdk.NewInt(1000),
					))
			}),
		),
		When("there are staking related transactions",
			Before(func() {
				oldMintInfo, newMintInfo, totalStaked, totalUnStaked, totalFee, err = s.mintInfoWithStakingTxs()
				Expect(err).To(BeNil())

				burnedFee = sdk.MustNewDecFromStr("0.5").MulInt(totalFee).TruncateInt()

				expNextInflation = s.NextInflationRate(oldMintInfo.Inflation, oldMintInfo.BondedRatio)
				expNextBlockProvision = expNextInflation.
					MulInt(oldMintInfo.CirculatingSupply).
					QuoInt(sdk.NewInt(int64(oldMintInfo.Params.InflationParameters.BlocksPerYear))).
					TruncateInt()

				expStakingRewards, expFoundation, expCommunity = getAllProportions(expNextBlockProvision, oldMintInfo.Params.InflationDistribution)
				_ = expStakingRewards
			}),

			It("bonded must adjust", func() {
				Expect(newMintInfo.StakingSupply.Sub(oldMintInfo.StakingSupply)).To(
					common.EQ(totalStaked.Sub(totalUnStaked)),
				)
			}),

			It("circulatingSupply must be correct", func() {
				Expect(newMintInfo.CirculatingSupply.Sub(oldMintInfo.CirculatingSupply)).To(
					Equal(expNextBlockProvision.Sub(burnedFee)))
			}),

			It("inflation must be correct", func() {
				Expect(newMintInfo.Inflation).To(common.LTE(s.params.InflationParameters.InflationMax))
				Expect(newMintInfo.Inflation).To(common.GTE(s.params.InflationParameters.InflationMin))
				Expect(newMintInfo.Inflation).To(Equal(expNextInflation))
			}),

			It("totalMintedProvision must properly increase", func() {
				Expect(newMintInfo.TotalMintedProvision.Sub(oldMintInfo.TotalMintedProvision)).To(Equal(
					expNextBlockProvision))
			}),

			It("foundation balance must be correct", func() {
				Expect(newMintInfo.FoundationBalance.Sub(oldMintInfo.FoundationBalance)).To(Equal(
					expFoundation,
				))
			}),

			It("community balance must be correct", func() {
				// we allow a variant of 1000.
				Expect(newMintInfo.CommunityBalance.Sub(oldMintInfo.CommunityBalance).Sub(expCommunity).Abs()).
					To(common.LTE(
						sdk.NewInt(1000),
					))
			}),
		),
	)

	return root
}
