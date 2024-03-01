package config

import (
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

// Load loads the config file.
// The config file should be named "config.yaml" and should be in the root of the project.
// Optionally, you can pass in a list of additional paths to search for the config file.
func Load(configPaths ...string) (Configuration, error) {
	var c Configuration

	// Load the config file from disk
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	for _, configPath := range configPaths {
		fileExtension := filepath.Ext(configPath)
		if fileExtension == ".yml" || fileExtension == ".yaml" {
			dir, file := filepath.Split(configPath)
			viper.SetConfigName(strings.Replace(file, fileExtension, "", 1))
			viper.AddConfigPath(dir)
		} else {
			viper.AddConfigPath(configPath)
		}
	}
	if err := viper.ReadInConfig(); err != nil {
		return c, err
	}

	// Parse the config into a struct.
	if err := viper.Unmarshal(&c); err != nil {
		return c, err
	}

	return c, nil
}
