package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port int `yaml:"port"`
}

// LoadConfig 从 YAML 配置文件加载配置
func LoadConfig(filePath string) (*Config, error) {
	var config Config
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
