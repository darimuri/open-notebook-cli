package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	APIURL  string
	APIKey  string
	Output  string
}

func Load(configPath string) (*Config, error) {
	v := viper.New()

	v.SetEnvPrefix("OPEN_NOTEBOOK")
	v.BindEnv("api_url", "OPEN_NOTEBOOK_API_URL")
	v.BindEnv("api_key", "OPEN_NOTEBOOK_API_KEY")
	v.BindEnv("output", "OPEN_NOTEBOOK_OUTPUT")

	v.SetDefault("api_url", "http://localhost:8080")
	v.SetDefault("output", "table")

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("$HOME/.config/open-notebook")
		v.AddConfigPath("/etc/open-notebook")
	}

	err := v.ReadInConfig()
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return &Config{
		APIURL: v.GetString("api_url"),
		APIKey: v.GetString("api_key"),
		Output: v.GetString("output"),
	}, nil
}