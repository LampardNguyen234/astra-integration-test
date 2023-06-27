package mint

import (
	mintTypes "github.com/AstraProtocol/astra/v2/x/mint/types"
	"github.com/LampardNguyen234/astra-go-sdk/client"
	"github.com/LampardNguyen234/astra-integration-test/common/logger"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/common"
)

type MintSuite struct {
	*common.TestClient
	params *mintTypes.Params
}

func NewMintSuite(cfg *SuiteConfig, cc *client.CosmosClient, log logger.Logger) (common.ITestSuite, error) {
	if _, err := cfg.IsValid(); err != nil {
		return nil, err
	}
	tc := common.NewTestClient(cc, log.WithPrefix("Mint Suite"))

	params, err := tc.MintParams()
	if err != nil {
		return nil, err
	}

	return &MintSuite{TestClient: tc, params: params}, nil
}

func (s *MintSuite) Name() string {
	return "mint"
}
