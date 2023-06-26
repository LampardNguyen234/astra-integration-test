package config

import (
	"encoding/json"
	"fmt"
	testSuite "github.com/LampardNguyen234/astra-integration-test/test-suite"
	"io/ioutil"
	"path/filepath"
)

type Config struct {
	Logger    LoggerConfig          `json:"Logger"`
	TestSuite testSuite.SuiteConfig `json:"TestSuite" json:"TestSuite"`
}

func DefaultConfig() Config {
	return Config{
		Logger:    DefaultLoggerConfig(),
		TestSuite: testSuite.DefaultConfig(),
	}
}

func LoadConfig(configFile string) (*Config, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("readFile %v error: %v", configFile, err)
	}

	var res Config
	ext := filepath.Ext(configFile)
	switch ext {
	case ".json":
		err = json.Unmarshal(data, &res)
	default:
		err = fmt.Errorf("file extension `%v` not supported", ext)
	}

	if err != nil {
		return nil, err
	}

	res.LoadEnv()
	if _, err = res.IsValid(); err != nil {
		return nil, err
	}

	return &res, nil
}

func (cfg *Config) IsValid() (bool, error) {
	if _, err := cfg.Logger.IsValid(); err != nil {
		return false, err
	}
	if _, err := cfg.TestSuite.IsValid(); err != nil {
		return false, err
	}

	return true, nil
}

func (cfg *Config) LoadEnv() {}
