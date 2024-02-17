package config

import (
	"github.com/spf13/viper"
)

// Load loads the config file.
// The config file should be named "config.yaml" and should be in the root of the project.
func Load() (Configuration, error) {
	var c Configuration

	// Load the config file from disk
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return c, err
	}

	// Parse the config into a struct.
	if err := viper.Unmarshal(&c); err != nil {
		return c, err
	}

	return c, nil
}
