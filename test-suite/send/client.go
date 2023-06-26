package send

import (
	"github.com/LampardNguyen234/astra-go-sdk/client"
	"github.com/LampardNguyen234/astra-integration-test/common/logger"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/common"
)

type SendSuite struct {
	*common.TestClient
}

func NewSendSuite(cfg *SuiteConfig, cc *client.CosmosClient, log logger.Logger) (*SendSuite, error) {
	if _, err := cfg.IsValid(); err != nil {
		return nil, err
	}
	tc := common.NewTestClient(cc, log.WithPrefix("Send Suite"))
	tc.SetMasterKey(cfg.MasterKey)

	return &SendSuite{TestClient: tc}, nil
}

func (s *SendSuite) Name() string {
	return "send"
}
