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
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/evmos/evmos/v6/x/vesting/types"
	"time"
)

type clawBackVestingTestCase struct {
	name             string
	txParams         msg_params.TxParams
	funder           *account.KeyInfo
	dest             string
	vestingAccount   *account.KeyInfo
	expClawBackedAmt sdk.Int
	expErr           error
	before           func()
	after            func()
}

func (tc clawBackVestingTestCase) Msg() sdk.Msg {
	var msg sdk.Msg
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

func (tc clawBackVestingTestCase) Name() string {
	return fmt.Sprintf("[tc: %v]", tc.name)
}

func (s *VestingSuite) runClawBackVestingTest() {

	//operators := common.RandKeyInfos(10)
	funders := common.RandKeyInfos(10)
	vestingAccounts := common.RandKeyInfos(10)
	defaultDest := account.MustNewPrivateKeyFromString(s.GetMasterKey()).AccAddress().String()

	createVestingFunc := func(funder, recipient *account.KeyInfo,
		start time.Time,
		vestingPeriod vestingTypes.Periods,
		lockPeriod vestingTypes.Periods,
	) {
		params := msg_params.TxParams{
			PrivateKey: funder.PrivateKey,
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
		assert.NoError(err)
		s.TxShouldPass(tx.TxHash)
	}

	tcs := []clawBackVestingTestCase{
		{
			name: "when coins are unvested but not locked - should claw back unvested",
			txParams: msg_params.TxParams{
				PrivateKey: funders[0].PrivateKey,
			},
			funder:           funders[0],
			dest:             defaultDest,
			vestingAccount:   vestingAccounts[0],
			expClawBackedAmt: sdk.NewIntWithDecimal(1, 15),
			before: func() {
				s.FundAccount(funders[0].CosmosAddress, 1)
				createVestingFunc(
					funders[0],
					vestingAccounts[0],
					time.Now(),
					vestingTypes.Periods{
						{
							Length: 1000,
							Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 15))),
						},
					},
					nil,
				)
			},
			after: func() {
				s.Refund(funders[0].PrivateKey)
				s.Refund(vestingAccounts[0].PrivateKey)
			},
			expErr: nil,
		},
		{
			name: "when coins are partially vested - should claw back unvested only",
			txParams: msg_params.TxParams{
				PrivateKey: funders[1].PrivateKey,
			},
			funder:           funders[1],
			dest:             defaultDest,
			vestingAccount:   vestingAccounts[1],
			expClawBackedAmt: sdk.NewIntWithDecimal(1, 15),
			before: func() {
				s.FundAccount(funders[1].CosmosAddress, 1)
				createVestingFunc(
					funders[1],
					vestingAccounts[1],
					time.Unix(time.Now().Unix()-1500, 0),
					vestingTypes.Periods{
						{
							Length: 1000,
							Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 15))),
						},
						{
							Length: 1000,
							Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 15))),
						},
					},
					nil,
				)
			},
			after: func() {
				s.Refund(funders[1].PrivateKey)
				s.Refund(vestingAccounts[1].PrivateKey)
			},
			expErr: nil,
		},
		{
			name: "when coins are unvested and locked - should claw back all",
			txParams: msg_params.TxParams{
				PrivateKey: funders[2].PrivateKey,
			},
			funder:           funders[2],
			dest:             defaultDest,
			vestingAccount:   vestingAccounts[2],
			expClawBackedAmt: sdk.NewIntWithDecimal(1, 15),
			before: func() {
				s.FundAccount(funders[2].CosmosAddress, 1)
				createVestingFunc(
					funders[2],
					vestingAccounts[2],
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
			},
			after: func() {
				s.Refund(funders[2].PrivateKey)
				s.Refund(vestingAccounts[2].PrivateKey)
			},
			expErr: nil,
		},
		{
			name: "when coins are vested and locked - should claw back unvested only",
			txParams: msg_params.TxParams{
				PrivateKey: funders[3].PrivateKey,
			},
			funder:           funders[3],
			dest:             defaultDest,
			vestingAccount:   vestingAccounts[3],
			expClawBackedAmt: sdk.NewIntWithDecimal(2, 15),
			before: func() {
				s.FundAccount(funders[3].CosmosAddress, 1)
				createVestingFunc(
					funders[3],
					vestingAccounts[3],
					time.Unix(time.Now().Unix()-1500, 0),
					vestingTypes.Periods{
						{
							Length: 1000,
							Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(2, 15))),
						},
						{
							Length: 1000,
							Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(2, 15))),
						},
					},
					vestingTypes.Periods{
						{
							Length: 1000,
							Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(1, 15))),
						},
						{
							Length: 1000,
							Amount: sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntWithDecimal(3, 15))),
						},
					},
				)
			},
			after: func() {
				s.Refund(funders[3].PrivateKey)
				s.Refund(vestingAccounts[3].PrivateKey)
			},
			expErr: nil,
		},
		//{
		//	name: "should be able to claw back with MsgExec",
		//},
		//{
		//	name: "should not be able to claw back on behalf of other",
		//},
	}

	for _, tc := range tcs {
		s.processClawBackVestingTestCase(tc)
	}

	s.Log.Infof("CLAWBACK VESTING PASS")
}

func (s *VestingSuite) processClawBackVestingTestCase(tc clawBackVestingTestCase) {
	msg := tc.Name()
	if tc.before != nil {
		tc.before()
	}
	defer func() {
		if tc.after != nil {
			tc.after()
		}
	}()

	oldVestingBalance, err := s.GetVestingBalance(tc.vestingAccount.CosmosAddress)
	assert.NoError(err)
	oldDestBalance, err := s.Balance(tc.dest)
	assert.NoError(err)

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

		newVestingBalance, err := s.GetVestingBalance(tc.vestingAccount.CosmosAddress)
		assert.NoError(err, msg)
		newDestBalance, err := s.Balance(tc.dest)
		assert.NoError(err, msg)

		// sanity checks
		assert.Equal(newVestingBalance.Unvested.AmountOf(sdkCommon.BaseDenom),
			sdk.ZeroInt(), msg)
		assert.Equal(newVestingBalance.Vested.AmountOf(sdkCommon.BaseDenom),
			oldVestingBalance.Vested.AmountOf(sdkCommon.BaseDenom), msg)
		assert.Equal(newVestingBalance.Locked.AmountOf(sdkCommon.BaseDenom),
			sdk.ZeroInt(), msg)
		assert.Equal(newVestingBalance.Unlocked.AmountOf(sdkCommon.BaseDenom),
			oldVestingBalance.Unlocked.AmountOf(sdkCommon.BaseDenom), msg)
		assert.Equal(oldVestingBalance.Total.Sub(newVestingBalance.Total).AmountOf(sdkCommon.BaseDenom), tc.expClawBackedAmt, msg)
		assert.Equal(newDestBalance.Total.Sub(oldDestBalance.Total), tc.expClawBackedAmt, msg)

		s.Log.Debugf("%v PASSED", msg)
	}
}
