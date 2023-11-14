package feeburn

import (
	"context"
	"fmt"
	"github.com/LampardNguyen234/astra-go-sdk/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
	"time"
)

func (s *FeeburnSuite) feeBurnInfo(withTx bool) (oldFeeBurn, newFeeBurn, totalFee sdk.Int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	oldFeeBurn, newFeeBurn, totalFee = sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt()
	sleepDuration := 200 * time.Millisecond
	for {
		select {
		case <-ctx.Done():
			err = fmt.Errorf("timeout")
			return
		default:
			var height *big.Int
			height, err = s.LatestBlockHeight()
			if err != nil {
				s.Log.Debugf("failed to retrieve block height: %v", err)
				time.Sleep(sleepDuration)
			}

			oldFeeBurn, err = s.TotalFeeBurn()
			if err != nil {
				s.Log.Debugf("failed to retrieve old feeburn: %v", err)
				time.Sleep(sleepDuration)
				continue
			}

			if withTx {
				s.RandomTxs()
			}
			s.WaitUntilBlock(height.Int64() + 1)

			newFeeBurn, err = s.TotalFeeBurn()
			if err != nil {
				s.Log.Debugf("failed to retrieve new feeburn: %v", err)
				time.Sleep(sleepDuration)
				continue
			}

			txs, err1 := s.BlockTxsByHeight(ctx, height.Add(height, big.NewInt(1)))
			if err1 != nil {
				s.Log.Debugf("failed to get txs for block %v: %v", height.Int64()+1, err1)
				time.Sleep(sleepDuration)
				continue
			}
			if !withTx && len(txs) > 0 {
				time.Sleep(sleepDuration)
				continue
			}
			if withTx && len(txs) == 0 {
				time.Sleep(sleepDuration)
				continue
			}

			for _, tx := range txs {
				txFee, _ := tx.GetTx().(sdk.FeeTx)
				if txFee == nil {
					continue
				}
				fmt.Println(totalFee)
				totalFee = totalFee.Add(common.ParseAmount(txFee.GetFee()))
			}
			err = nil
			return
		}
	}
}
