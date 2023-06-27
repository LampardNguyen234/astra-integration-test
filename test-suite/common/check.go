package common

import (
	"context"
	"strings"
	"time"
)

func (c *TestClient) TxShouldPass(txHash string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			c.Log.Panicf("failed to check txHash %v: TIMED-OUT", txHash)
		default:
			time.Sleep(2 * time.Second)
			resp, _ := c.TxByHash(txHash)
			if resp != nil {
				if resp.Height == 0 {
					continue
				}
				if resp.Code == 0 {
					c.WaitUntilBlock(resp.Height + 1)
					return
				} else {
					c.Log.Panicf("tx %v failed with code: %v, error: %v", resp.Code, resp.RawLog)
				}
			}
		}
	}
}

func (c *TestClient) TxShouldFailedWithError(txHash string, errMsg string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			c.Log.Panicf("failed to retrieve txHash %v: TIMED-OUT", txHash)
		default:
			time.Sleep(2 * time.Second)
			resp, _ := c.TxByHash(txHash)
			if resp != nil {
				if resp.Height == 0 {
					continue
				}
				if resp.Code == 0 {
					c.Log.Panicf("expect tx %v to fail with error: %v", txHash, errMsg)
				} else if !strings.Contains(resp.RawLog, errMsg) {
					c.Log.Panicf("expect tx %v to fail with error %v, got %v", errMsg, resp.RawLog)
				} else {
					return
				}
			}
		}
	}
}
