package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Server ServerConfig
	Log    LogConfig
}

type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type LogConfig struct {
	LogFilePath string `toml:"logfilepath"`
	MaxLogSize  int    `toml:"maxlogsize"`
}

// LoadConfig 从 TOML 配置文件加载配置
func LoadConfig(filePath string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
