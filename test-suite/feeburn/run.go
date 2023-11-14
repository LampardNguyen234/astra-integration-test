package feeburn

import (
	. "github.com/LampardNguyen234/astra-integration-test/framework"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func (s *FeeburnSuite) RunTest() {
	s.Start()

	tc := s.registerTests()
	tc.Run()
	tc.Report()

	s.Finished()
}

func (s *FeeburnSuite) registerTests() ITestNode {
	RegisterFailHandler(ginkgo.Fail)

	var oldBurnedFee, newBurnedFee sdk.Int
	var totalFee sdk.Int
	var err error

	root := Describe(s.Name(),
		When("there are no transactions",
			Before(func() {
				oldBurnedFee, newBurnedFee, totalFee, err = s.feeBurnInfo(false)
				Expect(err).To(BeNil())
			}),
			It("no fee should be burned",
				func() {
					Expect(totalFee.String()).To(Equal(sdk.ZeroInt().String()))
					Expect(newBurnedFee).To(Equal(oldBurnedFee))
				}),
		),
		When("there are transactions",
			Before(func() {
				oldBurnedFee, newBurnedFee, totalFee, err = s.feeBurnInfo(true)
				Expect(err).To(BeNil())
			}),

			It("50% fee must be burned", func() {
				Expect(newBurnedFee).To(Equal(
					oldBurnedFee.Add(sdk.NewDecFromInt(totalFee).Mul(sdk.MustNewDecFromStr("0.5")).TruncateInt()),
				))
			}),
		),
	)

	return root
}
