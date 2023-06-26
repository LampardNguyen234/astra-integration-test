package send

import (
	"fmt"
	"github.com/LampardNguyen234/astra-go-sdk/account"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/common"
)

type SuiteConfig struct {
	common.BaseConfig
	MasterKey string `json:"MasterKey"`
}

func (cfg *SuiteConfig) IsValid() (bool, error) {
	if !cfg.Enabled {
		return true, nil
	}
	if _, err := account.NewPrivateKeyFromString(cfg.MasterKey); err != nil {
		return false, fmt.Errorf("invalid `masterKey`: %v", err)
	}

	return true, nil
}

func DefaultConfig() SuiteConfig {
	return SuiteConfig{
		BaseConfig: common.DefaultBaseConfig(),
		MasterKey:  "",
	}
}
