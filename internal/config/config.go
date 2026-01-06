package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`

	Data struct {
		Dir     string `yaml:"dir"`
		AOFFile string `yaml:"aof_file"`
	} `yaml:"data"`

	Engine struct {
		CleanupIntervalSec int `yaml:"cleanup_interval_sec"`
	} `yaml:"engine"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, cfg)
	return cfg, err
}
