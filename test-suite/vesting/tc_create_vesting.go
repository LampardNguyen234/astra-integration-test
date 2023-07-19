package vesting

import (
	"fmt"
	"github.com/LampardNguyen234/astra-go-sdk/account"
	"github.com/LampardNguyen234/astra-go-sdk/client/msg_params"
	sdkCommon "github.com/LampardNguyen234/astra-go-sdk/common"
	"github.com/LampardNguyen234/astra-integration-test/common"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/assert"
	sdk "github.com/cosmos/cosmos-sdk/types"
	vestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/evmos/evmos/v6/x/vesting/types"
	"time"
)

type createVestingTestCase struct {
	name           string
	txParams       msg_params.TxParams
	funder         *account.KeyInfo
	recipient      *account.KeyInfo
	startTime      time.Time
	vestingPeriods vestingTypes.Periods
	lockupPeriods  vestingTypes.Periods
	before         func()
	after          func()
	expErr         error
}

func (tc createVestingTestCase) Msg() sdk.Msg {
	funder := tc.txParams.MustGetPrivateKey().AccAddress()
	if tc.funder != nil {
		funder = account.MustParseCosmosAddress(tc.funder.CosmosAddress)
	}

	return types.NewMsgCreateClawbackVestingAccount(
		funder,
		account.MustParseCosmosAddress(tc.recipient.CosmosAddress),
		tc.startTime,
		tc.lockupPeriods,
		tc.vestingPeriods,
	)
}

func (tc createVestingTestCase) Name() string {
	return fmt.Sprintf("[tc: %v]", tc.name)
}

func (s *VestingSuite) runCreateVestingTests() {
	funders := []*account.KeyInfo{
		common.RandKeyInfo(),
		common.RandKeyInfo(),
		common.RandKeyInfo(),
	}
	vestingAccounts := []*account.KeyInfo{
		common.RandKeyInfo(),
		common.RandKeyInfo(),
	}
	regularAccounts := []*account.KeyInfo{
		common.RandKeyInfo(),
	}
	tcs := []createVestingTestCase{
		{
			name: "should be able to create a vesting account with empty lockupPeriods",
			txParams: msg_params.TxParams{
				PrivateKey: funders[0].PrivateKey,
			},
			recipient:      common.RandKeyInfo(),
			startTime:      time.Now(),
			vestingPeriods: randPeriods(0.1),
			lockupPeriods:  nil,
			before: func() {
				s.FundAccount(funders[0].CosmosAddress, 0.2)
			},
			after: func() {
				s.Refund(funders[0].PrivateKey)
			},
			expErr: nil,
		},
		{
			name: "should be able to create a vesting account with non-empty lockupPeriods",
			txParams: msg_params.TxParams{
				PrivateKey: funders[0].PrivateKey,
			},
			recipient:      common.RandKeyInfo(),
			startTime:      time.Now(),
			vestingPeriods: randPeriods(0.1),
			lockupPeriods:  randPeriods(0.1),
			before: func() {
				s.FundAccount(funders[0].CosmosAddress, 0.2)
			},
			after: func() {
				s.Refund(funders[0].PrivateKey)
			},
			expErr: nil,
		},
		{
			name: "should be able to create a vesting account with non-empty lockupPeriods",
			txParams: msg_params.TxParams{
				PrivateKey: funders[0].PrivateKey,
			},
			recipient:      common.RandKeyInfo(),
			startTime:      time.Now(),
			vestingPeriods: randPeriods(0.1),
			lockupPeriods:  randPeriods(0.1),
			before: func() {
				s.FundAccount(funders[0].CosmosAddress, 0.2)
			},
			after: func() {
				s.Refund(funders[0].PrivateKey)
			},
			expErr: nil,
		},
		{
			name: "should fail to create a vesting account with lockupPeriods.TotalAmount() != vestingPeriods.TotalAmount()",
			txParams: msg_params.TxParams{
				PrivateKey: funders[0].PrivateKey,
			},
			recipient:      common.RandKeyInfo(),
			startTime:      time.Now(),
			vestingPeriods: randPeriods(0.1),
			lockupPeriods:  randPeriods(0.2),
			before: func() {
				s.FundAccount(funders[0].CosmosAddress, 0.2)
			},
			after: func() {
				s.Refund(funders[0].PrivateKey)
			},
			expErr: fmt.Errorf("vesting and lockup schedules must have same total coins"),
		},
		{
			name: "should fail to create a vesting account with zero vesting length",
			txParams: msg_params.TxParams{
				PrivateKey: funders[0].PrivateKey,
			},
			recipient: common.RandKeyInfo(),
			startTime: time.Now(),
			vestingPeriods: vestingTypes.Periods{
				{
					Length: 0,
					Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 18))),
				},
			},
			lockupPeriods: nil,
			before: func() {
				s.FundAccount(funders[0].CosmosAddress, 2)
			},
			after: func() {
				s.Refund(funders[0].PrivateKey)
			},
			expErr: fmt.Errorf("length must be greater than 0"),
		},
		{
			name: "should fail to create a vesting account with negative vesting length",
			txParams: msg_params.TxParams{
				PrivateKey: funders[0].PrivateKey,
			},
			recipient: common.RandKeyInfo(),
			startTime: time.Now(),
			vestingPeriods: vestingTypes.Periods{
				{
					Length: -1,
					Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 18))),
				},
			},
			lockupPeriods: nil,
			before: func() {
				s.FundAccount(funders[0].CosmosAddress, 2)
			},
			after: func() {
				s.Refund(funders[0].PrivateKey)
			},
			expErr: fmt.Errorf("length must be greater than 0"),
		},
		{
			name: "should fail to create a vesting account with negative locking length",
			txParams: msg_params.TxParams{
				PrivateKey: funders[0].PrivateKey,
			},
			recipient:      common.RandKeyInfo(),
			startTime:      time.Now(),
			vestingPeriods: randPeriods(1),
			lockupPeriods: vestingTypes.Periods{
				{
					Length: -1,
					Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 18))),
				},
			},
			before: func() {
				s.FundAccount(funders[0].CosmosAddress, 2)
			},
			after: func() {
				s.Refund(funders[0].PrivateKey)
			},
			expErr: fmt.Errorf("length must be greater than 0"),
		},
		{
			name: "should fail to create a vesting account with zero locking length",
			txParams: msg_params.TxParams{
				PrivateKey: funders[0].PrivateKey,
			},
			recipient:      common.RandKeyInfo(),
			startTime:      time.Now(),
			vestingPeriods: randPeriods(1),
			lockupPeriods: vestingTypes.Periods{
				{
					Length: 0,
					Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 18))),
				},
			},
			before: func() {
				s.FundAccount(funders[0].CosmosAddress, 2)
			},
			after: func() {
				s.Refund(funders[0].PrivateKey)
			},
			expErr: fmt.Errorf("length must be greater than 0"),
		},
		{
			name: "should fail to create a vesting account on behalf of others (invalid pubKey)",
			txParams: msg_params.TxParams{
				PrivateKey: funders[1].PrivateKey,
			},
			funder:         funders[2],
			recipient:      common.RandKeyInfo(),
			startTime:      time.Now(),
			vestingPeriods: randPeriods(0.1),
			lockupPeriods:  nil,
			before: func() {
				s.FundAccount(funders[1].CosmosAddress, 1)
				s.FundAccount(funders[2].CosmosAddress, 1)
			},
			after: func() {
				s.Refund(funders[1].PrivateKey)
				s.Refund(funders[2].PrivateKey)
			},
			expErr: fmt.Errorf("pubKey does not match signer address"),
		},
		{
			name: "should fail to create a vesting account upon an existing regular account",
			txParams: msg_params.TxParams{
				PrivateKey: funders[0].PrivateKey,
			},
			recipient:      regularAccounts[0],
			startTime:      time.Now(),
			vestingPeriods: randPeriods(0.1),
			lockupPeriods:  nil,
			before: func() {
				s.FundAccount(funders[0].CosmosAddress, 1)
				s.FundAccount(regularAccounts[0].CosmosAddress, 1)
			},
			after: func() {
				s.Refund(funders[0].PrivateKey)
				s.Refund(regularAccounts[0].PrivateKey)
			},
			expErr: fmt.Errorf("account %s must be a clawback vesting account", regularAccounts[0].CosmosAddress),
		},
		{
			name: "should fail to create a vesting account upon an existing vesting account",
			txParams: msg_params.TxParams{
				PrivateKey: funders[0].PrivateKey,
			},
			recipient:      vestingAccounts[0],
			startTime:      time.Now(),
			vestingPeriods: randPeriods(0.1),
			lockupPeriods:  nil,
			before: func() {
				s.FundAccount(funders[0].CosmosAddress, 0.2)
				s.FundVesting(vestingAccounts[0].CosmosAddress, 0.1, 0)
			},
			after: func() {
				s.Refund(funders[0].PrivateKey)
				s.Refund(vestingAccounts[0].PrivateKey)
			},
			expErr: fmt.Errorf("account %s already exists",
				vestingAccounts[0].CosmosAddress),
		},
		//{
		//	name: "should fail to create a vesting account with negative amount",
		//	txParams: msg_params.TxParams{
		//		PrivateKey: s.GetMasterKey(),
		//	},
		//	vestingAccount: common.RandKeyInfo(),
		//	startTime: time.Now(),
		//	vestingPeriods: vestingTypes.Periods{
		//		{
		//			Length: 1,
		//			Amount: sdk.Coins{
		//				{
		//					Denom:  sdkCommon.BaseDenom,
		//					Amount: sdk.NewIntWithDecimal(-1, 18),
		//				},
		//			},
		//		},
		//		{
		//			Length: 1,
		//			Amount: sdk.Coins{
		//				{
		//					Denom:  sdkCommon.BaseDenom,
		//					Amount: sdk.NewIntWithDecimal(-1, 18),
		//				},
		//			},
		//		},
		//		{
		//			Length: 1,
		//			Amount: sdk.Coins{
		//				{
		//					Denom:  sdkCommon.BaseDenom,
		//					Amount: sdk.NewIntWithDecimal(2, 18),
		//				},
		//			},
		//		},
		//	},
		//	lockupPeriods: nil,
		//	expErr:        fmt.Errorf("ádasdasda"),
		//},
	}

	for _, tc := range tcs {
		s.processCreateVestingTestCase(tc)
	}

	s.Log.Infof("CREATE VESTING PASS")
}

func (s *VestingSuite) processCreateVestingTestCase(tc createVestingTestCase) {
	msg := tc.Name()
	if tc.before != nil {
		tc.before()
	}
	defer func() {
		if tc.after != nil {
			tc.after()
		}
	}()

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
			assert.ErrorContains(err, tc.expErr.Error(), msg)
		}
		s.Log.Debugf("%v PASSED", msg)
	} else {
		assert.NoError(err, msg)
		s.TxShouldPass(resp.TxHash)

		va, err := s.GetVestingAccount(tc.recipient.CosmosAddress)
		assert.NoError(err, msg)
		balances, err := s.GetVestingBalance(tc.recipient.CosmosAddress)
		assert.NoError(err, msg)

		// sanity checks
		assert.Equal(account.MustParseCosmosAddress(va.FunderAddress).String(),
			funder.String(), msg)
		assert.Equal(va.StartTime.Unix(), tc.startTime.Unix(), msg)
		assert.Equal(va.StartTime.Unix(), tc.startTime.Unix(), msg)
		assert.Equal(va.VestingPeriods, tc.vestingPeriods, msg)
		assert.Equal(
			sdkCommon.ParseAmount(va.OriginalVesting),
			sdkCommon.ParseAmount(tc.vestingPeriods.TotalAmount()),
			msg,
		)

		assert.Compare(sdkCommon.ParseAmount(balances.Total),
			sdkCommon.ParseAmount(tc.vestingPeriods.TotalAmount()),
			assert.OpEQ,
		)

		if tc.lockupPeriods != nil {
			assert.Equal(va.LockupPeriods, tc.lockupPeriods, msg)
			assert.Compare(
				sdkCommon.ParseAmount(balances.Total),
				sdkCommon.ParseAmount(tc.lockupPeriods.TotalAmount()),
				assert.OpEQ,
			)
		} else {
			assert.Equal(va.LockupPeriods.TotalLength(), 0, msg)
			assert.Equal(
				sdkCommon.ParseAmount(va.LockupPeriods.TotalAmount()),
				sdkCommon.ParseAmount(tc.vestingPeriods.TotalAmount()),
				msg,
			)
		}

		s.Log.Debugf("%v PASSED", msg)
	}
}
