package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	StravaClientID       string `required:"true" split_words:"true"`
	StravaClientSecret   string `required:"true" split_words:"true"`
	FolderPath           string `required:"true" split_words:"true"`
	RefreshTokenFileName string `required:"true" split_words:"true" default:"refresh_token.json"`
}

func LoadConfig() (*Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
