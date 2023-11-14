package test_suite

import (
	"fmt"
	"github.com/LampardNguyen234/astra-go-sdk/account"
	"github.com/LampardNguyen234/astra-go-sdk/client"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/feeburn"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/mint"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/send"
	"github.com/LampardNguyen234/astra-integration-test/test-suite/vesting"
)

type SuiteConfig struct {
	MasterKey    string                    `json:"MasterKey,omitempty"`
	SdkConfig    client.CosmosClientConfig `json:"SdkConfig"`
	SendSuite    *send.SuiteConfig         `json:"SendSuite"`
	MintSuite    *mint.SuiteConfig         `json:"MintSuite"`
	VestingSuite *vesting.SuiteConfig      `json:"VestingSuite"`
	FeeBurnSuite *feeburn.SuiteConfig      `json:"FeeBurnSuite"`
}

func (cfg *SuiteConfig) IsValid() (bool, error) {
	if cfg.MasterKey != "" {
		if _, err := account.NewPrivateKeyFromString(cfg.MasterKey); err != nil {
			return false, fmt.Errorf("invalid master key: %v", err)
		}
	}
	if _, err := cfg.SdkConfig.IsValid(); err != nil {
		return false, fmt.Errorf("invalid sdk config: %v", err)
	}

	if _, err := cfg.SendSuite.IsValid(); err != nil {
		return false, fmt.Errorf("invalid send config: %v", err)
	}

	if _, err := cfg.VestingSuite.IsValid(); err != nil {
		return false, fmt.Errorf("invalid vesting config: %v", err)
	}

	if _, err := cfg.FeeBurnSuite.IsValid(); err != nil {
		return false, fmt.Errorf("invalid feeburn config: %v", err)
	}

	return true, nil
}

func DefaultConfig() SuiteConfig {
	ss := send.DefaultConfig()
	ms := mint.DefaultConfig()
	vs := vesting.DefaultConfig()
	fb := feeburn.DefaultConfig()
	return SuiteConfig{
		MasterKey:    "",
		SdkConfig:    client.DefaultTestnetConfig(),
		SendSuite:    &ss,
		MintSuite:    &ms,
		VestingSuite: &vs,
		FeeBurnSuite: &fb,
	}
}
