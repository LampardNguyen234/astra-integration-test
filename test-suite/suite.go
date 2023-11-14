package test_suite

import (
	"fmt"
	"github.com/LampardNguyen234/astra-go-sdk/client"
	"github.com/LampardNguyen234/astra-integration-test/common/logger"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/common"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/feeburn"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/mint"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/send"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/vesting"
)

type TestSuite struct {
	*client.CosmosClient
	log logger.Logger

	suites []common.ITestSuite
}

func NewTestSuite(cfg SuiteConfig, log logger.Logger) (*TestSuite, error) {
	if _, err := cfg.IsValid(); err != nil {
		return nil, err
	}
	c, err := client.NewCosmosClient(cfg.SdkConfig)
	if err != nil {
		return nil, err
	}

	suites := make([]common.ITestSuite, 0)
	if cfg.SendSuite.Enabled {
		if cfg.MasterKey != "" {
			cfg.SendSuite.MasterKey = cfg.MasterKey
		}
		ss, err := send.NewSendSuite(cfg.SendSuite, c, log)
		if err != nil {
			return nil, fmt.Errorf("invalid `send`: %v", err)
		}
		suites = append(suites, ss)
	}
	if cfg.MintSuite.Enabled {
		ms, err := mint.NewMintSuite(cfg.MintSuite, c, log)
		if err != nil {
			return nil, fmt.Errorf("invalid `mint`: %v", err)
		}
		if cfg.MasterKey != "" {
			ms.SetMasterKey(cfg.MasterKey)
		}
		suites = append(suites, ms)
	}
	if cfg.VestingSuite.Enabled {
		if cfg.MasterKey != "" {
			cfg.VestingSuite.MasterKey = cfg.MasterKey
		}
		ss, err := vesting.NewVestingSuite(cfg.VestingSuite, c, log)
		if err != nil {
			return nil, fmt.Errorf("invalid `vesting`: %v", err)
		}
		suites = append(suites, ss)
	}
	if cfg.FeeBurnSuite.Enabled {
		fb, err := feeburn.NewFeeburnSuite(cfg.FeeBurnSuite, c, log)
		if err != nil {
			return nil, fmt.Errorf("invalid `feeburn`: %v", err)
		}
		if cfg.MasterKey != "" {
			fb.SetMasterKey(cfg.MasterKey)
		}
		suites = append(suites, fb)
	}

	return &TestSuite{
		CosmosClient: c,
		suites:       suites,
		log:          log.WithPrefix("Main"),
	}, nil
}

func (s *TestSuite) RunTest() {
	for _, suite := range s.suites {
		suite.RunTest()
	}
}
