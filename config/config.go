package config

import (
	"os"

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
	if !FileExists(filePath) {
		// 楔入配置文件
		err := DefaultConfig().WriteConfig(filePath)
		if err != nil {
			return nil, err
		}
		return DefaultConfig(), nil
	}

	var config Config
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// 写入配置文件
func (c *Config) WriteConfig(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(c)
}

// 检测文件是否存在
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

/*
[server]
host = "127.0.0.1"
port = 8080

[mmdb]
mmdbpath = "/data/ip/db"
asndbpath = "/data/ip/db/asn.mmdb"
countrydbpath = "/data/ip/db/country.mmdb"
ipinfoKey = ""
updateFreq = 24 # hours

[log]
logfilepath = "/data/ip/log/ip.log"
maxlogsize = 5 # MB
*/

func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Log: LogConfig{
			LogFilePath: "./log/ip.log",
			MaxLogSize:  5, // MB
		},
		Mmdb: MmdbConfig{
			MmDBPath:      "./data",
			ASNDBPath:     "./data/asn.mmdb",
			CountryDBPath: "./data/country.mmdb",
			IPinfoKey:     "",
			UpdateFreq:    24, // hours
		},
	}
}
