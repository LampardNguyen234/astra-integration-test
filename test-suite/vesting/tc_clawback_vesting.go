package vesting

import (
	"github.com/LampardNguyen234/astra-go-sdk/account"
	"github.com/LampardNguyen234/astra-go-sdk/client/msg_params"
	sdkCommon "github.com/LampardNguyen234/astra-go-sdk/common"
	"github.com/LampardNguyen234/astra-integration-test/common"
	. "github.com/LampardNguyen234/astra-integration-test/framework"
	sdk "github.com/cosmos/cosmos-sdk/types"
	vestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/evmos/evmos/v6/x/vesting/types"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

type clawBackVestingTestCase struct {
	txParams         msg_params.TxParams
	funder           *account.KeyInfo
	operator         *account.KeyInfo
	dest             string
	vestingAccount   *account.KeyInfo
	expClawBackedAmt sdk.Int
	expErr           error
}

func (tc clawBackVestingTestCase) Msg() sdk.Msg {
	var msg sdk.Msg
	dest := tc.dest
	if dest == "" {
		dest = tc.funder.CosmosAddress
	}
	msg = types.NewMsgClawback(
		account.MustParseCosmosAddress(tc.funder.CosmosAddress),
		account.MustParseCosmosAddress(tc.vestingAccount.CosmosAddress),
		account.MustParseCosmosAddress(tc.dest),
	)
	if tc.txParams.PrivateKey != tc.funder.PrivateKey {
		tmp := authz.NewMsgExec(
			account.MustNewPrivateKeyFromString(tc.txParams.PrivateKey).AccAddress(),
			[]sdk.Msg{
				msg,
			},
		)
		msg = &tmp
	}

	return msg
}

func (s *VestingSuite) testClawBackVesting() ITestNode {
	RegisterFailHandler(ginkgo.Fail)
	var funder, operator, vestingAccount *account.KeyInfo
	var txParams msg_params.TxParams
	defaultDest := account.MustNewPrivateKeyFromString(s.GetMasterKey()).AccAddress().String()

	defaultBeforeEach := BeforeEach(func() {
		funder = common.RandKeyInfo()
		vestingAccount = common.RandKeyInfo()
		operator = common.RandKeyInfo()
		txParams = msg_params.TxParams{
			PrivateKey: funder.PrivateKey,
			GasLimit:   300000,
		}
		_ = operator

		s.FundAccount(funder.CosmosAddress, 1)
	})
	defaultAfterEach := AfterEach(func() {
		s.Refund(funder.PrivateKey)
		s.Refund(operator.PrivateKey)
	})

	createVestingFunc := func(funder, recipient *account.KeyInfo,
		start time.Time,
		vestingPeriod vestingTypes.Periods,
		lockPeriod vestingTypes.Periods,
	) {
		params := msg_params.TxParams{
			PrivateKey: funder.PrivateKey,
			GasLimit:   500000,
		}

		tx, err := s.BuildAndSendTx(
			params,
			types.NewMsgCreateClawbackVestingAccount(
				account.MustParseCosmosAddress(funder.CosmosAddress),
				account.MustParseCosmosAddress(recipient.CosmosAddress),
				start,
				vestingPeriod,
				vestingPeriod,
			),
		)
		Expect(err).To(BeNil())
		s.TxShouldPass(tx.TxHash)
	}

	root := Describe("clawback vesting",
		When("coins are unvested",
			defaultBeforeEach,
			defaultAfterEach,
			Context("coins are unlocked",
				Before(func() {
					createVestingFunc(
						funder,
						vestingAccount,
						time.Now(),
						vestingTypes.Periods{
							{
								Length: 1000,
								Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 15))),
							},
						},
						nil,
					)
				}),
				It("should clawback all", func() {
					s.processClawBackVestingTestCase(clawBackVestingTestCase{
						txParams:         txParams,
						funder:           funder,
						dest:             defaultDest,
						vestingAccount:   vestingAccount,
						expClawBackedAmt: sdk.NewIntWithDecimal(1, 15),
						expErr:           nil,
					})
				}),
			),

			Context("coins are locked",
				Before(func() {
					createVestingFunc(
						funder,
						vestingAccount,
						time.Now(),
						vestingTypes.Periods{
							{
								Length: 1000,
								Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 15))),
							},
						},
						vestingTypes.Periods{
							{
								Length: 100,
								Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 15))),
							},
						},
					)
				}),
				It("should clawback all", func() {
					s.processClawBackVestingTestCase(clawBackVestingTestCase{
						txParams:         txParams,
						funder:           funder,
						dest:             defaultDest,
						vestingAccount:   vestingAccount,
						expClawBackedAmt: sdk.NewIntWithDecimal(1, 15),
						expErr:           nil,
					})
				}),
			),
		),

		When("coins are vested",
			Context("coins are locked",
				BeforeEach(func() {
					defaultBeforeEach()
					createVestingFunc(
						funder,
						vestingAccount,
						time.Unix(time.Now().Unix()-10, 0),
						vestingTypes.Periods{
							{
								Length: 1,
								Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 15))),
							},
						},
						vestingTypes.Periods{
							{
								Length: 100,
								Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 15))),
							},
						},
					)
				}),
				defaultAfterEach,
				It("should clawback nothing", func() {
					s.processClawBackVestingTestCase(clawBackVestingTestCase{
						txParams:         txParams,
						funder:           funder,
						dest:             defaultDest,
						vestingAccount:   vestingAccount,
						expClawBackedAmt: sdk.ZeroInt(),
						expErr:           nil,
					})
				}),
			),
			Context("coins are unlocked",
				BeforeEach(func() {
					defaultBeforeEach()
					createVestingFunc(
						funder,
						vestingAccount,
						time.Unix(time.Now().Unix()-10, 0),
						vestingTypes.Periods{
							{
								Length: 1,
								Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 15))),
							},
						},
						vestingTypes.Periods{
							{
								Length: 1,
								Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 15))),
							},
						},
					)
				}),
				defaultAfterEach,
				It("should clawback nothing", func() {
					s.processClawBackVestingTestCase(clawBackVestingTestCase{
						txParams:         txParams,
						funder:           funder,
						dest:             defaultDest,
						vestingAccount:   vestingAccount,
						expClawBackedAmt: sdk.ZeroInt(),
						expErr:           nil,
					})
				}),
			),
		),
	)

	return root
}

func (s *VestingSuite) processClawBackVestingTestCase(tc clawBackVestingTestCase) {
	oldVestingBalance, err := s.GetVestingBalance(tc.vestingAccount.CosmosAddress)
	Expect(err).To(BeNil())
	oldDestBalance, err := s.Balance(tc.dest)
	Expect(err).To(BeNil())

	tc.txParams.GasAdjustment = 1.5
	resp, err := s.BuildAndSendTx(
		tc.txParams,
		tc.Msg(),
	)
	if tc.expErr != nil {
		if err == nil {
			s.TxShouldFailWithError(resp.TxHash, tc.expErr.Error())
		} else {
			Expect(err.Error()).To(ContainSubstring(tc.expErr.Error()))
		}
	} else {
		Expect(err).To(BeNil())
		s.TxShouldPass(resp.TxHash)

		newVestingBalance, err := s.GetVestingBalance(tc.vestingAccount.CosmosAddress)
		Expect(err).To(BeNil())
		newDestBalance, err := s.Balance(tc.dest)
		Expect(err).To(BeNil())

		// sanity checks
		Expect(newVestingBalance.Unvested.AmountOf(sdkCommon.BaseDenom)).To(Equal(
			sdk.ZeroInt()))
		Expect(newVestingBalance.Vested.AmountOf(sdkCommon.BaseDenom)).To(Equal(
			oldVestingBalance.Vested.AmountOf(sdkCommon.BaseDenom)))
		Expect(newVestingBalance.Locked.AmountOf(sdkCommon.BaseDenom)).To(Equal(
			sdk.ZeroInt()))
		Expect(newVestingBalance.Unlocked.AmountOf(sdkCommon.BaseDenom)).To(Equal(
			oldVestingBalance.Unlocked.AmountOf(sdkCommon.BaseDenom)))
		s.Log.Debugf("%v, %v, %v, %v",
			oldVestingBalance.Total.AmountOf(sdkCommon.BaseDenom),
			newVestingBalance.Total.AmountOf(sdkCommon.BaseDenom),
			oldVestingBalance.Total.Sub(newVestingBalance.Total).AmountOf(sdkCommon.BaseDenom),
			tc.expClawBackedAmt,
		)
		Expect(oldVestingBalance.Total.Sub(newVestingBalance.Total).AmountOf(sdkCommon.BaseDenom)).To(Equal(
			tc.expClawBackedAmt))
		Expect(newDestBalance.Total.Sub(oldDestBalance.Total).String()).To(Equal(tc.expClawBackedAmt.String()))
	}
}
