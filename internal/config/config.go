package config

import (
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

	// Only read config file if path is explicitly provided
	if configPath == "" {
		// Try to find config in default paths, ignore error if not found
		_ = v.ReadInConfig()
	} else {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	return &Config{
		APIURL: v.GetString("api_url"),
		APIKey: v.GetString("api_key"),
		Output: v.GetString("output"),
	}, nil
}