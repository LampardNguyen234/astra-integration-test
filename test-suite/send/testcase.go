package send

import (
	"fmt"
	"github.com/LampardNguyen234/astra-go-sdk/account"
	"github.com/LampardNguyen234/astra-go-sdk/client/msg_params"
	sdkCommon "github.com/LampardNguyen234/astra-go-sdk/common"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/assert"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type sendTestCase struct {
	name      string
	txParams  msg_params.TxParams
	recipient *account.KeyInfo
	amt       float64
	prefunded float64
	mall      func()
	expErr    error
}

func (s *SendSuite) runTestCase(tc sendTestCase) {
	msg := fmt.Sprintf("[tc: `%v`]", tc.name)
	s.FundAccount(tc.txParams.MustGetPrivateKey().AccAddress().String(), tc.prefunded)
	if tc.mall != nil {
		tc.mall()
	}
	defer s.ClawBack(tc.txParams.PrivateKey)

	// balance before sending
	rcptBal, err := s.CosmosClient.Balance(tc.recipient.CosmosAddress)
	assert.NoError(err, msg)

	resp, err := s.CosmosClient.BuildAndSendTx(
		tc.txParams,
		bankTypes.NewMsgSend(
			account.MustParseCosmosAddress(tc.txParams.MustGetPrivateKey().AccAddress().String()),
			account.MustParseCosmosAddress(tc.recipient.CosmosAddress),
			sdk.NewCoins(sdk.NewCoin(sdkCommon.BaseDenom, sdk.NewIntFromBigInt(sdkCommon.Float64ToBigInt(tc.amt)))),
		),
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
		s.BalanceCompare(tc.recipient.CosmosAddress,
			rcptBal.Total.Add(sdk.NewIntFromBigInt(sdkCommon.Float64ToBigInt(tc.amt))),
			assert.OpEQ,
		)
		s.Log.Debugf("%v PASSED", msg)
	}
}
