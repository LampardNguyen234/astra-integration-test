package mint

import (
	mintTypes "github.com/AstraProtocol/astra/v2/x/mint/types"
	"github.com/LampardNguyen234/astra-go-sdk/client"
	"github.com/LampardNguyen234/astra-integration-test/common/logger"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (s *MintSuite) NextInflationRate(inflation, bondedRatio sdk.Dec) sdk.Dec {
	// (1 - bondedRatio/GoalBonded) * InflationRateChange
	inflationRateChangePerYear := sdk.OneDec().
		Sub(bondedRatio.Quo(s.params.InflationParameters.GoalBonded)).
		Mul(s.params.InflationParameters.InflationRateChange)

	inflationRateChange := inflationRateChangePerYear.QuoInt(sdk.NewIntFromUint64(s.params.InflationParameters.BlocksPerYear))
	// adjust the new annual inflation for this next cycle
	inflation = inflation.Add(inflationRateChange) // note inflationRateChange may be negative
	if inflation.GT(s.params.InflationParameters.InflationMax) {
		inflation = s.params.InflationParameters.InflationMax
	}
	if inflation.LT(s.params.InflationParameters.InflationMin) {
		inflation = s.params.InflationParameters.InflationMin
	}

	return inflation
}
