package mint

import (
	"github.com/LampardNguyen234/astra-integration-test/test-suite/common"
)

type SuiteConfig struct {
	common.BaseConfig
}

func (cfg *SuiteConfig) IsValid() (bool, error) {
	if !cfg.Enabled {
		return true, nil
	}

	return true, nil
}

func DefaultConfig() SuiteConfig {
	return SuiteConfig{
		BaseConfig: common.DefaultBaseConfig(),
	}
}
