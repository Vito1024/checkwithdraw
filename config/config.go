package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CheckWithdraw `yaml:"check_withdraw"`
	OkLinkConfig  `yaml:"oklink"`
	UnisatConfig  `yaml:"unisat"`
}

func New(configFilepath string) *Config {
	var config Config

	bs, err := os.ReadFile(configFilepath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(bs, &config)
	if err != nil {
		panic(err)
	}

	return &config
}

type OkLinkConfig struct {
	Host string `yaml:"host"`
	Key  string `yaml:"key"`
}

type UnisatConfig struct {
	Host string `yaml:"host"`
}

type CheckWithdraw struct {
	ExcludedAddresses []string `yaml:"excluded_addresses"`
}
