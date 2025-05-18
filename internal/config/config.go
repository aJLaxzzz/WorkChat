package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	ConfigPathEnvKey  = "CONFIG_PATH"
	DefaultConfigPath = "./config.yaml"
	TemplatesDirPath  = "./templates"
)

type Config struct {
	CookiesSecretKey string `yaml:"cookies_secret_key"`
	EncryptionKey    string `yaml:"encryption_key"`
	Server           struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	DB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		Host     string `yaml:"host"`
	} `yaml:"db"`
}

func NewConfig() (*Config, error) {
	configPath := DefaultConfigPath
	if envConfigPath := os.Getenv(ConfigPathEnvKey); envConfigPath != "" {
		configPath = envConfigPath
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}

	var config Config
	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("error decoding the config: %v", err)
	}

	return &config, nil
}
