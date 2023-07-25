package common

import (
	"context"
	"time"
)

func (c *TestClient) WaitForBlock(numBlocks int64) {
	if numBlocks > 200 {
		c.Log.Panic("too many blocks, maximum supported: 200")
	}
	if numBlocks <= 0 {
		return
	}

	current, err := c.LatestBlockHeight()
	if err != nil {
		c.Log.Panic(err.Error())
	}

	c.WaitUntilBlock(current.Int64() + numBlocks)
}

func (c *TestClient) WaitUntilBlock(blk int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			c.Log.Panicf("failed to wait until block %v: TIMED-OUT", blk)
		default:
			time.Sleep(10 * time.Millisecond)
			resp, _ := c.LatestBlockHeight()
			if resp != nil {
				if resp.Int64() >= blk {
					return
				}
			}
		}
	}
}
