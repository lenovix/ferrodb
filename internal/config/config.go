package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
	} `yaml:"server"`

	Data struct {
		Dir     string `yaml:"dir"`
		AOFFile string `yaml:"aof_file"`
	} `yaml:"data"`

	Engine struct {
		DBCount            int `yaml:"db_count"`
		CleanupIntervalSec int `yaml:"cleanup_interval_sec"`
	} `yaml:"engine"`
}

func Load(path string) (*Config, error) {
	cfg := defaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	cfg.applyDefaults()
	return cfg, nil
}

func defaultConfig() *Config {
	cfg := &Config{}

	cfg.Server.Address = ":6380"
	cfg.Server.Password = ""

	cfg.Data.Dir = "data"
	cfg.Data.AOFFile = "ferrodb.aof"

	cfg.Engine.DBCount = 16
	cfg.Engine.CleanupIntervalSec = 1

	return cfg
}

func (c *Config) applyDefaults() {
	if c.Server.Address == "" {
		c.Server.Address = ":6380"
	}

	if c.Data.Dir == "" {
		c.Data.Dir = "data"
	}

	if c.Data.AOFFile == "" {
		c.Data.AOFFile = "ferrodb.aof"
	}

	if c.Engine.DBCount <= 0 {
		c.Engine.DBCount = 16
	}

	if c.Engine.CleanupIntervalSec <= 0 {
		c.Engine.CleanupIntervalSec = 1
	}
}

// Helper: full path to AOF
func (c *Config) AOFPath() string {
	return filepath.Join(c.Data.Dir, c.Data.AOFFile)
}
