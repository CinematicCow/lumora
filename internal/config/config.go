package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const ConfigFileName = "config"

type Config struct {
	DefaultDB string            `mapstructure:"default_db"`
	DBPaths   map[string]string `mapstructure:"db_paths"`
}

func InitConfig() (*Config, error) {
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "lumora")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	if _, err := os.Stat(filepath.Join(configDir, ConfigFileName+".yaml")); os.IsNotExist(err) {
		cfg := &Config{
			DefaultDB: "default",
			DBPaths:   make(map[string]string),
		}

		if len(cfg.DBPaths) == 0 {
			cfg.DBPaths["default"] = filepath.Join(configDir, "default")
		}
		viper.Set("default_db", cfg.DefaultDB)
		viper.Set("db_paths", cfg.DBPaths)
		err = viper.SafeWriteConfig()
		if err != nil {
			return nil, err
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Save() error {
	viper.Set("default_db", c.DefaultDB)
	viper.Set("db_paths", c.DBPaths)
	return viper.WriteConfig()
}

func (c *Config) AddDB(name, path string) {
	c.DBPaths[name] = path
}

func (c *Config) GetDBPath(name string) (string, bool) {
	path, exists := c.DBPaths[name]
	return path, exists
}

func (c *Config) SetDefaultDB(name string) {
	c.DefaultDB = name
}
