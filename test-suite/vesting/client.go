package vesting

import (
	"github.com/LampardNguyen234/astra-go-sdk/client"
	"github.com/LampardNguyen234/astra-go-sdk/client/msg_params"
	sdkCommon "github.com/LampardNguyen234/astra-go-sdk/common"
	"github.com/LampardNguyen234/astra-integration-test/common/logger"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/assert"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/common"
	"time"
)

type VestingSuite struct {
	*common.TestClient
}

func NewVestingSuite(cfg *SuiteConfig, cc *client.CosmosClient, log logger.Logger) (*VestingSuite, error) {
	if _, err := cfg.IsValid(); err != nil {
		return nil, err
	}
	tc := common.NewTestClient(cc, log.WithPrefix("Vesting Suite"))
	tc.SetMasterKey(cfg.MasterKey)

	return &VestingSuite{TestClient: tc}, nil
}

func (s *VestingSuite) Name() string {
	return "VestingModule"
}

func (s *VestingSuite) FundVesting(recipient string, amount float64, duration int64) {
	if duration == 0 {
		duration = 600
	}

	resp, err := s.TxCreateVesting(msg_params.TxCreateVestingParams{
		TxParams: msg_params.TxParams{
			PrivateKey: s.GetMasterKey(),
			Memo:       "Fund vesting from master",
		},
		ToAddr:          recipient,
		Start:           time.Now(),
		VestingDuration: duration,
		Amount:          sdkCommon.Float64ToBigInt(amount),
		VestingLength:   10,
	})
	assert.NoError(err)
	s.TxShouldPass(resp.TxHash)
}
