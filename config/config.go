package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Server ServerConfig
	Log    LogConfig
	Mmdb   MmdbConfig
}

type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type LogConfig struct {
	LogFilePath string `toml:"logfilepath"`
	MaxLogSize  int    `toml:"maxlogsize"`
}

type MmdbConfig struct {
	MmDBPath      string `toml:"mmdbpath"`
	ASNDBPath     string `toml:"asndbpath"`
	CountryDBPath string `toml:"countrydbpath"`
	IPinfoKey     string `toml:"ipinfoKey"`
	UpdateFreq    int    `toml:"updateFreq"`
}

// LoadConfig 从 TOML 配置文件加载配置
func LoadConfig(filePath string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
