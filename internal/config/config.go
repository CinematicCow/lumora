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
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(os.Getenv("HOME"), ".config", "lumora"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// config not found, create one
			config := &Config{
				DefaultDB: "",
				DBPaths:   make(map[string]string),
			}
			if err := config.Save(); err != nil {
				return nil, err
			}
			return config, nil
		} else {
			return nil, err
		}
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
