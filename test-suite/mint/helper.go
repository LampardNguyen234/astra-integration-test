package mint

import (
	"context"
	"fmt"
	"github.com/LampardNguyen234/astra-go-sdk/client"
	"github.com/LampardNguyen234/astra-go-sdk/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"math/big"
	"time"
)

func (s *MintSuite) mintInfoWithNoStakingTxs() (old, new *client.ProvisionInfo, totalFee sdk.Int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	totalFee = sdk.ZeroInt()
	for {
		select {
		case <-ctx.Done():
			err = fmt.Errorf("time-out with err: %v", ctx.Err())
			return
		default:
			old, err = s.MintInfo()
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}

			s.RandomTxs()

			s.WaitUntilBlock(old.Height + 1)

			new, err = s.MintInfo()
			if err != nil || new.Height != old.Height+1 {
				time.Sleep(1 * time.Second)
				continue
			}

			txs, _ := s.BlockTxsByHeight(ctx, big.NewInt(new.Height))
			if len(txs) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}
			for _, tx := range txs {
				txFee, _ := tx.GetTx().(sdk.FeeTx)
				if txFee == nil {
					continue
				}
				totalFee = totalFee.Add(common.ParseAmount(txFee.GetFee()))
				failed := false
				for _, msg := range txFee.GetMsgs() {
					switch sdk.MsgTypeURL(msg) {
					case sdk.MsgTypeURL(&stakingTypes.MsgDelegate{}),
						sdk.MsgTypeURL(&stakingTypes.MsgUndelegate{}),
						sdk.MsgTypeURL(&stakingTypes.MsgBeginRedelegate{}):
						sdk.MsgTypeURL(&distrTypes.MsgWithdrawDelegatorReward{})
						sdk.MsgTypeURL(&distrTypes.MsgWithdrawValidatorCommission{})
						failed = true
						break
					}
				}
				if failed {
					continue
				}
			}

			return
		}
	}
}

func (s *MintSuite) mintInfoWithStakingTxs() (old, new *client.ProvisionInfo, totalStaked, totalUnStaked, totalFee sdk.Int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	totalFee = sdk.ZeroInt()
	totalStaked = sdk.ZeroInt()
	totalUnStaked = sdk.ZeroInt()
	for {
		select {
		case <-ctx.Done():
			err = fmt.Errorf("time-out with err: %v", ctx.Err())
			return
		default:
			old, err = s.MintInfo()
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}

			s.RandomTxs(
				sdk.MsgTypeURL(&stakingTypes.MsgDelegate{}),
				sdk.MsgTypeURL(&stakingTypes.MsgUndelegate{}),
			)

			s.WaitUntilBlock(old.Height + 1)

			new, err = s.MintInfo()
			if err != nil || new.Height != old.Height+1 {
				time.Sleep(1 * time.Second)
				continue
			}

			txs, _ := s.BlockTxsByHeight(ctx, big.NewInt(new.Height))
			if len(txs) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}
			for _, tx := range txs {
				txFee, _ := tx.GetTx().(sdk.FeeTx)
				if txFee == nil {
					continue
				}
				totalFee = totalFee.Add(common.ParseAmount(txFee.GetFee()))
				failed := true
				for _, msg := range txFee.GetMsgs() {
					switch sdk.MsgTypeURL(msg) {
					case sdk.MsgTypeURL(&stakingTypes.MsgDelegate{}):
						totalStaked = totalStaked.Add(msg.(*stakingTypes.MsgDelegate).Amount.Amount)
						failed = false
					case sdk.MsgTypeURL(&stakingTypes.MsgUndelegate{}):
						totalUnStaked = totalUnStaked.Add(msg.(*stakingTypes.MsgUndelegate).Amount.Amount)
						failed = false
					case sdk.MsgTypeURL(&stakingTypes.MsgBeginRedelegate{}):
						sdk.MsgTypeURL(&distrTypes.MsgWithdrawDelegatorReward{})
						sdk.MsgTypeURL(&distrTypes.MsgWithdrawValidatorCommission{})
						failed = false
					}
				}
				if failed {
					continue
				}
			}

			return
		}
	}
}
