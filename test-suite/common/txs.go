package common

import (
	"github.com/LampardNguyen234/astra-go-sdk/account"
	"github.com/LampardNguyen234/astra-go-sdk/client/msg_params"
	"github.com/LampardNguyen234/astra-go-sdk/common"
	repoCommon "github.com/LampardNguyen234/astra-integration-test/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"time"
)

func (c *TestClient) RandomTxs(msgTypes ...string) {
	var randFromList = func(msgTypes ...string) string {
		return msgTypes[int(repoCommon.RandUint64()%uint64(len(msgTypes)))]
	}

	var msg sdk.Msg
	msgType := sdk.MsgTypeURL(&types.MsgSend{})
	if len(msgTypes) > 0 {
		msgType = randFromList(msgTypes...)
	}
	switch msgType {
	case sdk.MsgTypeURL(&types.MsgSend{}):
		msg = types.NewMsgSend(
			c.Address(),
			account.MustParseCosmosAddress(repoCommon.RandKeyInfo().CosmosAddress),
			sdk.NewCoins(sdk.NewCoin(common.BaseDenom, sdk.NewIntWithDecimal(
				int64(1+repoCommon.RandUint64()%10), 15))),
		)
	case sdk.MsgTypeURL(&authz.MsgGrant{}):
		msg, _ = authz.NewMsgGrant(
			c.Address(),
			account.MustParseCosmosAddress(repoCommon.RandKeyInfo().CosmosAddress),
			authz.NewGenericAuthorization(sdk.MsgTypeURL(&authz.MsgGrant{})),
			time.Unix(time.Now().Unix()+10, 10),
		)
	case sdk.MsgTypeURL(&stakingTypes.MsgDelegate{}):
		msg = stakingTypes.NewMsgDelegate(
			c.Address(),
			account.MustParseCosmosValidatorAddress(c.MustRandActiveValidator().OperatorAddress),
			sdk.NewCoin(common.BaseDenom, sdk.NewIntWithDecimal(int64(1+repoCommon.RandUint64()%10), 10)),
		)
	case sdk.MsgTypeURL(&stakingTypes.MsgUndelegate{}):
		allDelegations, err := c.DelegationDetail(c.Address().String())
		if err == nil {
			var validator sdk.ValAddress
			var amt sdk.Int
			for val, detail := range allDelegations {
				validator = account.MustParseCosmosValidatorAddress(val)
				amt = detail.Amount.Quo(sdk.NewInt(2))
			}
			if !validator.Empty() {
				msg = stakingTypes.NewMsgUndelegate(
					c.Address(),
					validator,
					sdk.NewCoin(common.BaseDenom, amt),
				)
			}
		}
	case sdk.MsgTypeURL(&stakingTypes.MsgBeginRedelegate{}):
		allDelegations, err := c.DelegationDetail(c.Address().String())
		if err == nil {
			var validator sdk.ValAddress
			var amt sdk.Int
			for val, detail := range allDelegations {
				validator = account.MustParseCosmosValidatorAddress(val)
				amt = detail.Amount
			}
			if !validator.Empty() {
				msg = stakingTypes.NewMsgBeginRedelegate(
					c.Address(),
					validator,
					account.MustParseCosmosValidatorAddress(c.MustRandActiveValidator().OperatorAddress),
					sdk.NewCoin(common.BaseDenom, amt),
				)
			}
		}
	}

	if msg == nil {
		return
	}
	_, _ = c.BuildAndSendTx(
		msg_params.TxParams{
			PrivateKey: c.GetMasterKey(),
			Memo:       "Testing Random Tx",
		},
		msg,
	)
}
