package mint

import (
	"fmt"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/common"
)

type SuiteConfig struct {
	common.BaseConfig
	NumTests int `json:"NumTests"`
}

func (cfg *SuiteConfig) IsValid() (bool, error) {
	if !cfg.Enabled {
		return true, nil
	}
	if cfg.NumTests <= 0 {
		return false, fmt.Errorf("invalid numTests %v", cfg.NumTests)
	}

	return true, nil
}

func DefaultConfig() SuiteConfig {
	return SuiteConfig{
		BaseConfig: common.DefaultBaseConfig(),
		NumTests:   10,
	}
}
