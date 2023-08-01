package vesting

import (
	"fmt"
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

type createVestingTestCase struct {
	txParams       msg_params.TxParams
	funder         *account.KeyInfo
	operator       *account.KeyInfo
	recipient      *account.KeyInfo
	startTime      time.Time
	vestingPeriods vestingTypes.Periods
	lockupPeriods  vestingTypes.Periods
	expErr         error
}

func (tc createVestingTestCase) Msg() sdk.Msg {
	funder := tc.txParams.MustGetPrivateKey().AccAddress()
	if tc.funder != nil {
		funder = account.MustParseCosmosAddress(tc.funder.CosmosAddress)
	}

	var msg sdk.Msg
	msg = types.NewMsgCreateClawbackVestingAccount(
		funder,
		account.MustParseCosmosAddress(tc.recipient.CosmosAddress),
		tc.startTime,
		tc.lockupPeriods,
		tc.vestingPeriods,
	)
	if tc.operator != nil && tc.operator.PrivateKey != tc.funder.PrivateKey {
		tmp := authz.NewMsgExec(account.MustParseCosmosAddress(tc.operator.CosmosAddress),
			[]sdk.Msg{msg},
		)
		msg = &tmp
	}

	return msg
}

func (s *VestingSuite) testCreateVesting() ITestNode {
	RegisterFailHandler(ginkgo.Fail)
	var funder, operator, to *account.KeyInfo
	var txParams msg_params.TxParams

	defaultBeforeEach := BeforeEach(func() {
		funder = common.RandKeyInfo()
		to = common.RandKeyInfo()
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

	root := Describe("create vesting",
		defaultBeforeEach,
		defaultAfterEach,

		It("should be able to create a vesting account with empty lockupPeriods", func() {
			s.processCreateVestingTestCase(createVestingTestCase{
				txParams:       txParams,
				funder:         funder,
				recipient:      to,
				startTime:      time.Now(),
				vestingPeriods: randPeriods(0.1),
				lockupPeriods:  nil,
			})
		}),

		It("should be able to create a vesting account with non-empty lockupPeriods", func() {
			s.processCreateVestingTestCase(createVestingTestCase{
				txParams:       txParams,
				funder:         funder,
				recipient:      to,
				startTime:      time.Now(),
				vestingPeriods: randPeriods(0.1),
				lockupPeriods:  randPeriods(0.1),
				expErr:         nil,
			})
		}),

		It("should be able to create a vesting account with startTime in the past",
			After(func() {
				s.Refund(to.PrivateKey)
			}),
			func() {
				s.processCreateVestingTestCase(createVestingTestCase{
					txParams:       txParams,
					funder:         funder,
					recipient:      to,
					startTime:      time.Now().Add(time.Duration(-(1 + common.RandInt64()%100))),
					vestingPeriods: randPeriods(0.1),
					lockupPeriods:  randPeriods(0.1),
				})
			}),

		It("should be able to create a vesting account with startTime in the future",
			func() {
				s.processCreateVestingTestCase(createVestingTestCase{
					txParams:       txParams,
					funder:         funder,
					recipient:      to,
					startTime:      time.Now().Add(time.Duration(1 + common.RandInt64()%100)),
					vestingPeriods: randPeriods(0.1),
					lockupPeriods:  randPeriods(0.1),
				})
			}),

		It("should not be able to create a vesting account on behalf of granter (unsupported)",
			Before(func() {
				s.FundAccount(operator.CosmosAddress, 0.1)
				resp, err := s.CosmosClient.TxGrantAuthorization(
					msg_params.TxGrantParams{
						TxParams:    txParams,
						Grantee:     operator.CosmosAddress,
						ExpiredTime: time.Now().Add(100),
					},
					authz.NewGenericAuthorization(sdk.MsgTypeURL(&types.MsgCreateClawbackVestingAccount{})),
				)
				Expect(err).To(BeNil())
				s.TxShouldPass(resp.TxHash)
			}),
			func() {
				s.processCreateVestingTestCase(createVestingTestCase{
					txParams: msg_params.TxParams{
						PrivateKey: operator.PrivateKey,
					},
					funder:         funder,
					operator:       operator,
					recipient:      to,
					startTime:      time.Now().Add(time.Duration(1 + common.RandInt64()%100)),
					vestingPeriods: randPeriods(0.1),
					lockupPeriods:  randPeriods(0.1),
					expErr:         fmt.Errorf("authorization not found"),
				})
			}),

		It("should fail to create a vesting account with lockupPeriods.TotalAmount() != vestingPeriods.TotalAmount()", func() {
			s.processCreateVestingTestCase(createVestingTestCase{
				txParams:       txParams,
				funder:         funder,
				recipient:      to,
				startTime:      time.Now(),
				vestingPeriods: randPeriods(0.1),
				lockupPeriods:  randPeriods(0.2),
				expErr:         fmt.Errorf("vesting and lockup schedules must have same total coins"),
			})
		}),

		It("should fail to create a vesting account with zero vesting length", func() {
			s.processCreateVestingTestCase(createVestingTestCase{
				txParams:  txParams,
				funder:    funder,
				recipient: to,
				startTime: time.Now(),
				vestingPeriods: vestingTypes.Periods{
					{
						Length: 10,
						Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 16))),
					},
					{
						Length: 0,
						Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(3, 16))),
					},
				},
				lockupPeriods: nil,
				expErr:        fmt.Errorf("length must be greater than 0"),
			})
		}),

		It("should fail to create a vesting account with negative vesting length", func() {
			s.processCreateVestingTestCase(createVestingTestCase{
				txParams:  txParams,
				funder:    funder,
				recipient: to,
				startTime: time.Now(),
				vestingPeriods: vestingTypes.Periods{
					{
						Length: 10,
						Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 16))),
					},
					{
						Length: -1,
						Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(3, 16))),
					},
				},
				lockupPeriods: nil,
				expErr:        fmt.Errorf("length must be greater than 0"),
			})
		}),

		It("should fail to create a vesting account with zero locking length", func() {
			s.processCreateVestingTestCase(createVestingTestCase{
				txParams:       txParams,
				funder:         funder,
				recipient:      to,
				startTime:      time.Now(),
				vestingPeriods: randPeriods(0.1),
				lockupPeriods: vestingTypes.Periods{
					{
						Length: 0,
						Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 17))),
					},
				},
				expErr: fmt.Errorf("length must be greater than 0"),
			})
		}),

		It("should fail to create a vesting account with negative locking length", func() {
			s.processCreateVestingTestCase(createVestingTestCase{
				txParams:       txParams,
				funder:         funder,
				recipient:      to,
				startTime:      time.Now(),
				vestingPeriods: randPeriods(0.1),
				lockupPeriods: vestingTypes.Periods{
					{
						Length: -1,
						Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 17))),
					},
				},
				expErr: fmt.Errorf("length must be greater than 0"),
			})
		}),

		It("should fail to create a vesting account on behalf of others (invalid pubKey)", func() {
			fakeFunder := common.RandKeyInfo()
			s.FundAccount(fakeFunder.CosmosAddress, 1)
			s.processCreateVestingTestCase(createVestingTestCase{
				txParams: msg_params.TxParams{
					PrivateKey: fakeFunder.PrivateKey,
				},
				funder:         funder,
				recipient:      to,
				startTime:      time.Now(),
				vestingPeriods: randPeriods(0.1),
				lockupPeriods:  randPeriods(0.1),
				expErr:         fmt.Errorf("pubKey does not match signer address"),
			})
		}),

		It("should fail to create a vesting account upon an existing regular account", func() {
			regularAccount := common.RandKeyInfo()
			s.FundAccount(regularAccount.CosmosAddress, 1)
			s.processCreateVestingTestCase(createVestingTestCase{
				txParams:       txParams,
				funder:         funder,
				recipient:      regularAccount,
				startTime:      time.Now(),
				vestingPeriods: randPeriods(0.1),
				lockupPeriods:  randPeriods(0.1),
				expErr:         fmt.Errorf("account %s already exists", regularAccount.CosmosAddress),
			})
		}),

		It("should fail to create a vesting account upon an existing vesting account", func() {
			s.FundVesting(to.CosmosAddress, 1, 0)
			s.processCreateVestingTestCase(createVestingTestCase{
				txParams:       txParams,
				funder:         funder,
				recipient:      to,
				startTime:      time.Now(),
				vestingPeriods: randPeriods(0.1),
				lockupPeriods:  randPeriods(0.1),
				expErr:         fmt.Errorf("account %s already exists", to.CosmosAddress),
			})
		}),

		It("should fail to create a vesting account with negative amount"),
	)

	return root
}

func (s *VestingSuite) processCreateVestingTestCase(tc createVestingTestCase) {
	RegisterFailHandler(ginkgo.Fail)
	funder := tc.txParams.MustGetPrivateKey().AccAddress()
	if tc.funder != nil {
		funder = account.MustParseCosmosAddress(tc.funder.CosmosAddress)
	}

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

		va, err := s.GetVestingAccount(tc.recipient.CosmosAddress)
		Expect(err).To(BeNil())
		balances, err := s.GetVestingBalance(tc.recipient.CosmosAddress)
		Expect(err).To(BeNil())

		// sanity checks
		Expect(account.MustParseCosmosAddress(va.FunderAddress).String()).To(
			Equal(funder.String()),
		)
		Expect(va.StartTime.Unix()).To(Equal(tc.startTime.Unix()))
		Expect(va.StartTime.Unix()).To(Equal(tc.startTime.Unix()))
		Expect(va.VestingPeriods).To(Equal(tc.vestingPeriods))
		Expect(sdkCommon.ParseAmount(va.OriginalVesting)).To(
			Equal(sdkCommon.ParseAmount(tc.vestingPeriods.TotalAmount())),
		)
		Expect(sdkCommon.ParseAmount(balances.Total)).To(Equal(sdkCommon.ParseAmount(tc.vestingPeriods.TotalAmount())))
		if tc.lockupPeriods != nil {
			Expect(va.LockupPeriods).To(Equal(tc.lockupPeriods))
			Expect(sdkCommon.ParseAmount(balances.Total)).To(Equal(
				sdkCommon.ParseAmount(tc.lockupPeriods.TotalAmount()),
			))
		} else {
			Expect(va.LockupPeriods.TotalLength()).To(Equal(int64(0)))
			Expect(sdkCommon.ParseAmount(va.LockupPeriods.TotalAmount())).To(Equal(
				sdkCommon.ParseAmount(tc.vestingPeriods.TotalAmount())),
			)
		}
	}
}
