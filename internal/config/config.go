package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	APIURL    string
	APIKey    string
	Output    string
	Notebook  string
}

func Load(configPath string) (*Config, error) {
	v := viper.New()

	v.SetEnvPrefix("OPEN_NOTEBOOK")
	v.BindEnv("api_url", "API_URL")
	v.BindEnv("api_key", "API_KEY")
	v.BindEnv("output", "OUTPUT")

	v.SetDefault("api_url", "https://open-notebook.darimuri.me")
	v.SetDefault("output", "table")

	if configPath == "" {
		home, _ := os.UserHomeDir()
		v.AddConfigPath(home + "/.config/open-notebook")
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		_ = v.ReadInConfig()
	} else {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	return &Config{
		APIURL:   v.GetString("api_url"),
		APIKey:   v.GetString("api_key"),
		Output:   v.GetString("output"),
		Notebook: v.GetString("notebook"),
	}, nil
}