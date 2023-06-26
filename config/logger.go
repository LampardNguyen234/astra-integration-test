package config

type LoggerConfig struct {
	LogPath string `yaml:"log_path" toml:"log_path" xml:"log_path" json:"LogPath"`
	Level   int8   `yaml:"level" toml:"level" xml:"level" json:"Level"`
	Color   bool   `yaml:"color" toml:"color" xml:"color" json:"Color"`
}

func (cfg *LoggerConfig) IsValid() (bool, error) {
	return true, nil
}

func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		LogPath: "",
		Level:   0,
		Color:   false,
	}
}
