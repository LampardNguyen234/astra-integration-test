package feeburn

import (
	mintTypes "github.com/AstraProtocol/astra/v2/x/mint/types"
	"github.com/LampardNguyen234/astra-go-sdk/client"
	"github.com/LampardNguyen234/astra-integration-test/common/logger"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/common"
)

type FeeburnSuite struct {
	*common.TestClient
	params *mintTypes.Params
}

func NewFeeburnSuite(cfg *SuiteConfig, cc *client.CosmosClient, log logger.Logger) (common.ITestSuite, error) {
	if _, err := cfg.IsValid(); err != nil {
		return nil, err
	}
	tc := common.NewTestClient(cc, log.WithPrefix("Feeburn Suite"))

	params, err := tc.MintParams()
	if err != nil {
		return nil, err
	}

	return &FeeburnSuite{TestClient: tc, params: params}, nil
}

func (s *FeeburnSuite) Name() string {
	return "mint"
}
