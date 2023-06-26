package common

type BaseConfig struct {
	Enabled bool `json:"Enabled"`
}

func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		Enabled: true,
	}
}
